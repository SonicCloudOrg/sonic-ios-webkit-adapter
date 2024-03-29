/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-webkit-adapter/entity/WebKitProtocol"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/yezihack/e"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func initProtocolAdapter(adapter *Adapter, version string) *protocolAdapter {
	protocol := &protocolAdapter{
		adapter: adapter,
	}
	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		if major <= 8 {
			initIOS8(protocol)
			return protocol
		}
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		if major > 12 || major >= 12 && minor >= 2 {
			initIOS12(protocol)
			return protocol
		}
	}
	initIOS9(protocol)
	return protocol
}

type mapSelectorListFunc func(selectorList *WebKitProtocol.SelectorList)

type protocolAdapter struct {
	adapter                    *Adapter
	lastNodeId                 int64
	lastPageExecutionContextId int64
	styleMap                   map[string]interface{}
	lastScriptEval             interface{}
	screencast                 *screencastSession
	mapSelectorList            mapSelectorListFunc
}

func (p *protocolAdapter) init() {
	p.mapSelectorList = func(selectorList *WebKitProtocol.SelectorList) {

	}
	p.styleMap = make(map[string]interface{})

	p.adapter.AddToolMessageFilter("DOM.getDocument", p.onDomGetDocument)
	// CSS
	p.adapter.AddToolMessageFilter("CSS.setStyleTexts", p.onSetStyleTexts)
	p.adapter.AddToolMessageFilter("CSS.getMatchedStylesForNode", p.onGetMatchedStylesForNode)
	p.adapter.AddToolMessageFilter("CSS.getBackgroundColors", p.onGetBackgroundColors)
	p.adapter.AddToolMessageFilter("CSS.addRule", p.onAddRule)
	p.adapter.AddToolMessageFilter("CSS.getPlatformFontsForNode", p.onGetPlatformFontsForNode)

	p.adapter.AddWebkitMessageFilter("CSS.getMatchedStylesForNode", p.onGetMatchedStylesForNodeResult)
	// Page
	p.adapter.AddToolMessageFilter("Page.startScreencast", p.onStartScreencast)
	p.adapter.AddToolMessageFilter("Page.stopScreencast", p.onStopScreencast)
	p.adapter.AddToolMessageFilter("Page.screencastFrameAck", p.onScreencastFrameAck)
	p.adapter.AddToolMessageFilter("Page.getNavigationHistory", p.onGetNavigationHistory)
	p.adapter.AddToolMessageFilter("Page.setOverlayMessage", p.onPageSetOverlay)
	p.adapter.AddToolMessageFilter("Page.configureOverlay", p.onPageConfigureOverlay)
	// DOM
	p.adapter.AddToolMessageFilter("DOM.enable", p.onDomEnable)
	p.adapter.AddToolMessageFilter("DOM.setInspectMode", p.onSetInspectMode)
	p.adapter.AddToolMessageFilter("DOM.setInspectedNode", p.onDomSetInspectedNode)
	p.adapter.AddToolMessageFilter("DOM.pushNodesByBackendIdsToFrontend", p.onPushNodesByBackendIdsToFrontend)
	p.adapter.AddToolMessageFilter("DOM.getBoxModel", p.onGetBoxModel)
	p.adapter.AddToolMessageFilter("DOM.getNodeForLocation", p.onGetNodeForLocation)
	// DOMDebugger
	p.adapter.AddToolMessageFilter("DOMDebugger.getEventListeners", p.domDebuggerOnGetEventListeners)
	// Debugger
	p.adapter.AddToolMessageFilter("Debugger.canSetScriptSource", p.onCanSetScriptSource)
	p.adapter.AddToolMessageFilter("Debugger.setBlackboxPatterns", p.onSetBlackboxPatterns)
	p.adapter.AddToolMessageFilter("Debugger.setAsyncCallStackDepth", p.onSetAsyncCallStackDepth)
	p.adapter.AddToolMessageFilter("Debugger.enable", p.onDebuggerEnable)

	p.adapter.AddWebkitMessageFilter("Debugger.scriptParsed", p.onScriptParsed)
	// Emulation
	p.adapter.AddToolMessageFilter("Emulation.canEmulate", p.onCanEmulate)
	p.adapter.AddToolMessageFilter("Emulation.setTouchEmulationEnabled", p.onEmulationSetTouchEmulationEnabled)
	p.adapter.AddToolMessageFilter("Emulation.setScriptExecutionDisabled", p.onEmulationSetScriptExecutionDisabled)
	p.adapter.AddToolMessageFilter("Emulation.setEmulatedMedia", p.onEmulationSetEmulatedMedia)
	// Rendering
	p.adapter.AddToolMessageFilter("Rendering.setShowPaintRects", p.onRenderingSetShowPaintRects)
	// Input
	p.adapter.AddToolMessageFilter("Input.emulateTouchFromMouseEvent", p.onEmulateTouchFromMouseEvent)
	// Log
	p.adapter.AddToolMessageFilter("Log.clear", p.onLogClear)
	p.adapter.AddToolMessageFilter("Log.disable", p.onLogDisable)
	p.adapter.AddToolMessageFilter("Log.enable", p.onLogEnable)
	// Console
	p.adapter.AddWebkitMessageFilter("Console.messageAdded", p.onConsoleMessageAdded)
	// Network
	p.adapter.AddToolMessageFilter("Network.getCookies", p.onNetworkGetCookies)
	p.adapter.AddToolMessageFilter("Network.deleteCookie", p.onNetworkDeleteCookie)
	p.adapter.AddToolMessageFilter("Network.setMonitoringXHREnabled", p.onNetworkSetMonitoringXHREnabled)
	p.adapter.AddToolMessageFilter("Network.canEmulateNetworkConditions", p.onCanEmulateNetworkConditions)
	// Runtime
	p.adapter.AddToolMessageFilter("Runtime.compileScript", p.onRuntimeOnCompileScript)
	p.adapter.AddWebkitMessageFilter("Runtime.executionContextCreated", p.onExecutionContextCreated)
	p.adapter.AddWebkitMessageFilter("Runtime.evaluate", p.onEvaluate)
	p.adapter.AddWebkitMessageFilter("Runtime.getProperties", p.onRuntimeGetProperties)
	// Inspector
	p.adapter.AddToolMessageFilter("Inspector.inspect", p.onInspect)
}

func (p *protocolAdapter) defaultCallFunc(message []byte) {
	//log.Println(string(message))
}

func (p *protocolAdapter) onDomGetDocument(message []byte) []byte {
	p.enumerateStyleSheets(message)
	return message
}

func (p *protocolAdapter) onPageSetOverlay(message []byte) []byte {
	method := "Debugger.setOverlayMessage"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onPageConfigureOverlay(message []byte) []byte {
	return p.onPageSetOverlay(message)
}

func (p *protocolAdapter) onDomSetInspectedNode(message []byte) []byte {
	method := "Console.addInspectedNode"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onEmulationSetTouchEmulationEnabled(message []byte) []byte {
	method := "Page.setTouchEmulationEnabled"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onEmulationSetScriptExecutionDisabled(message []byte) []byte {
	method := "Page.setScriptExecutionDisabled"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onEmulationSetEmulatedMedia(message []byte) []byte {
	method := "Page.setEmulatedMedia"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onRenderingSetShowPaintRects(message []byte) []byte {
	method := "Page.setShowPaintRects"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onLogClear(message []byte) []byte {
	method := "Console.clearMessages"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onLogDisable(message []byte) []byte {
	method := "Console.disable"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onLogEnable(message []byte) []byte {
	method := "Console.enable"
	return ReplaceMethodNameAndOutputBinary(message, method)
}
func (p *protocolAdapter) onNetworkGetCookies(message []byte) []byte {
	method := "Page.getCookies"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onNetworkDeleteCookie(message []byte) []byte {
	method := "Page.deleteCookie"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onNetworkSetMonitoringXHREnabled(message []byte) []byte {
	method := "Console.setMonitoringXHREnabled"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *protocolAdapter) onGetMatchedStylesForNode(message []byte) []byte {
	p.lastNodeId = gjson.Get(string(message), "params.nodeId").Int()
	return message
}

func (p *protocolAdapter) onCanEmulate(message []byte) []byte {
	result := map[string]interface{}{
		"result": true,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onGetPlatformFontsForNode(message []byte) []byte {
	result := map[string]interface{}{
		"fonts": []string{},
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onGetBackgroundColors(message []byte) []byte {
	result := map[string]interface{}{
		"backgroundColors": []string{},
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onCanSetScriptSource(message []byte) []byte {
	result := map[string]interface{}{
		"result": false,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onSetBlackboxPatterns(message []byte) []byte {
	result := map[string]interface{}{}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onSetAsyncCallStackDepth(message []byte) []byte {
	result := map[string]interface{}{
		"result": true,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onDebuggerEnable(message []byte) []byte {
	p.adapter.CallTarget("Debugger.setBreakpointsActive", map[string]interface{}{
		"active": true,
	}, p.defaultCallFunc)
	return message
}

func (p *protocolAdapter) onExecutionContextCreated(message []byte) []byte {
	msg := string(message)
	var err error
	if gjson.Get(msg, "params").Exists() && gjson.Get(msg, "params.context").Exists() {
		if !gjson.Get(msg, "params.context.origin").Exists() {
			msg, err = sjson.Set(msg, "params.context.origin", gjson.Get(msg, "params.context.name").String())
			if err != nil {
				log.Println(e.Convert(err).ToStr())
			}
			if gjson.Get(msg, "params.context.isPageContext").Exists() {
				p.lastPageExecutionContextId = gjson.Get(msg, "params.context.id").Int()
			}
			if gjson.Get(msg, "params.context.frameId").Exists() {
				msg, err = sjson.Set(msg, "params.context.auxData", map[string]interface{}{
					"frameId":   gjson.Get(msg, "params.context.frameId").String(),
					"isDefault": true,
				})
				if err != nil {
					log.Println(e.Convert(err).ToStr())
				}
				if gjson.Get(msg, "params.context.frameId").Exists() {
					msg, err = sjson.Delete(msg, "params.context.frameId")
					if err != nil {
						log.Println(e.Convert(err).ToStr())
					}
				}
			}
		}
	}

	return []byte(msg)
}

func (p *protocolAdapter) onEvaluate(message []byte) []byte {
	msg := string(message)
	var err error
	result := gjson.Get(msg, "result")
	if result.Exists() && result.Get("wasThrown").Exists() {
		msg, err = sjson.Set(msg, "result.result.subtype", "error")
		if err != nil {
			return nil
		}
		msg, err = sjson.Set(msg, "result.exceptionDetails", map[string]interface{}{
			"text":     gjson.Get(msg, "result.result.description").Value(),
			"url":      "",
			"scriptId": p.lastScriptEval,
			"line":     1,
			"column":   0,
			"stack": map[string]interface{}{
				"callFrames": []map[string]interface{}{
					{
						"functionName": "",
						"scriptId":     p.lastScriptEval,
						"url":          "",
						"lineNumber":   1,
						"columnNumber": 1,
					},
				},
			},
		})
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
	} else if result.Exists() && result.Get("result").Exists() && result.Get("result.preview").Exists() {
		msg, err = sjson.Set(msg, "result.result.preview.description", gjson.Get(msg, "result.result.description").Value())
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		msg, err = sjson.Set(msg, "result.result.preview.type", "object")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
	}
	return []byte(msg)
}

func (p *protocolAdapter) onRuntimeOnCompileScript(message []byte) []byte {
	params := map[string]interface{}{
		"expression": gjson.Get(string(message), "params.expression").String(),
		"contextId":  gjson.Get(string(message), "params.executionContextId").Int(),
	}
	p.adapter.CallTarget("Runtime.evaluate", params, func(msg []byte) {
		result := map[string]interface{}{
			"scriptId":         nil,
			"exceptionDetails": nil,
		}
		p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	})
	return nil
}

func (p *protocolAdapter) onRuntimeGetProperties(message []byte) []byte {
	var newPropertyDescriptors []interface{}
	var err error
	msg := string(message)
	for _, node := range gjson.Get(msg, "result.result").Array() {
		isOwn := node.Get("isOwn")
		nativeGetter := node.Get("nativeGetter")
		if isOwn.Exists() || nativeGetter.Exists() {
			msg, err = sjson.Set(msg, isOwn.Path(string(message)), true)
			if err != nil {
				log.Println(e.Convert(err).ToStr())
			}
			newPropertyDescriptors = append(newPropertyDescriptors, gjson.Get(msg, node.Path(msg)).Value())
		}
	}
	msg, err = sjson.Set(msg, "result.result", newPropertyDescriptors)
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	return []byte(msg)
}

func (p *protocolAdapter) onScriptParsed(message []byte) []byte {
	p.lastScriptEval = gjson.Get(string(message), "params.scriptId")
	return message
}

func (p *protocolAdapter) onDomEnable(message []byte) []byte {
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), map[string]interface{}{})
	return nil
}

func (p *protocolAdapter) onSetInspectMode(message []byte) []byte {
	msg := string(message)
	var err error
	msg, err = sjson.Set(msg, "method", "DOM.setInspectModeEnabled")
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	msg, err = sjson.Set(msg, "params.enabled",
		gjson.Get(msg, "params.mode").String() == "searchForNode")
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	if gjson.Get(msg, "params.mode").Exists() {
		msg, err = sjson.Delete(msg, "params.mode")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
	}
	return []byte(msg)
}

func (p *protocolAdapter) onInspect(message []byte) []byte {
	msg := string(message)
	var err error
	msg, err = sjson.Set(msg, "method", "DOM.inspectNodeRequested")
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	msg, err = sjson.Set(msg, "params.backendNodeId", gjson.Get(msg, "params.object.objectId").Value())
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	if gjson.Get(msg, "params.hints").Exists() {
		msg, err = sjson.Delete(msg, "params.object")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
	}
	if gjson.Get(msg, "params.hints").Exists() {
		msg, err = sjson.Delete(msg, "params.hints")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
	}

	return []byte(msg)
}

func (p *protocolAdapter) domDebuggerOnGetEventListeners(message []byte) []byte {
	requestNodeParams := map[string]interface{}{
		"objectId": gjson.Get(string(message), "params.objectId").Value(),
	}
	p.adapter.CallTarget("DOM.requestNode", requestNodeParams, func(result []byte) {
		getEventListenersForNodeParams := map[string]interface{}{
			"nodeId":      gjson.Get(string(result), "result.nodeId").Value(),
			"objectGroup": "event-listeners-panel",
		}
		p.adapter.CallTarget("DOM.getEventListenersForNode", getEventListenersForNodeParams, func(msg []byte) {
			var getEventListenersForNodeResult = &WebKitProtocol.GetEventListenersForNodeResult{}
			err := json.Unmarshal(msg, getEventListenersForNodeResult)
			if err != nil {
				log.Println(e.Convert(err).ToStr())
			}
			listeners := getEventListenersForNodeResult.Listeners
			var mappedListeners []map[string]interface{}
			for _, listener := range listeners {
				mappedListeners = append(mappedListeners, map[string]interface{}{
					"type":       listener.Type,
					"useCapture": listener.UseCapture,
					"passive":    false, // iOS doesn't support this property, http://compatibility.remotedebug.org/DOM/Safari%20iOS%209.3/types/EventListener,
					"location":   listener.Location,
					"hander":     listener.HandlerName,
				})
			}
			mappedResult := map[string]interface{}{
				"listeners": mappedListeners,
			}
			p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), mappedResult)
		})
	})
	return nil
}

func (p *protocolAdapter) onPushNodesByBackendIdsToFrontend(message []byte) []byte {
	id := gjson.Get(string(message), "id").Int()
	var resultBackNodeIds []interface{}
	for _, backNode := range gjson.Get(string(message), "params.backendNodeIds").Array() {
		params := map[string]interface{}{
			"backendNodeId": backNode.Value(),
		}
		p.adapter.CallTarget("DOM.pushNodeByBackendIdToFrontend", params, func(msg []byte) {
			resultBackNodeIds = append(resultBackNodeIds, gjson.Get(string(msg), "nodeId").Value())
		})
	}
	result := map[string]interface{}{
		"nodeIds": resultBackNodeIds,
	}
	p.adapter.FireResultToTools(int(id), result)
	return nil
}

func (p *protocolAdapter) onGetBoxModel(message []byte) []byte {
	params := map[string]interface{}{
		"highlightConfig": map[string]interface{}{
			"showInfo":           true,
			"showRulers":         false,
			"showExtensionLines": false,
			"contentColor":       map[string]interface{}{"r": 111, "g": 168, "b": 220, "a": 0.66},
			"paddingColor":       map[string]interface{}{"r": 147, "g": 196, "b": 125, "a": 0.55},
			"borderColor":        map[string]interface{}{"r": 255, "g": 229, "b": 153, "a": 0.66},
			"marginColor":        map[string]interface{}{"r": 246, "g": 178, "b": 107, "a": 0.66},
			"eventTargetColor":   map[string]interface{}{"r": 255, "g": 196, "b": 196, "a": 0.66},
			"shapeColor":         map[string]interface{}{"r": 96, "g": 82, "b": 177, "a": 0.8},
			"shapeMarginColor":   map[string]interface{}{"r": 96, "g": 82, "b": 127, "a": 0.6},
			"displayAsMaterial":  true,
		},
		"nodeId": gjson.Get(string(message), "params.nodeId").Value(),
	}
	p.adapter.CallTarget("DOM.highlightNode", params, func(message []byte) {

	})
	return nil
}

func (p *protocolAdapter) onGetNodeForLocation(message []byte) []byte {
	evaluateParams := map[string]interface{}{
		"expression": fmt.Sprintf("document.elementFromPoint(%d,%d)", gjson.Get(string(message), "params.x").Int(), gjson.Get(string(message), "params.y").Int()),
	}
	p.adapter.CallTarget("Runtime.evaluate", evaluateParams, func(result []byte) {
		requestNodeParams := map[string]interface{}{
			"objectId": gjson.Get(string(result), "result.objectId").Value(),
		}
		p.adapter.CallTarget("DOM.requestNode", requestNodeParams, func(msg []byte) {
			resultParams := map[string]interface{}{
				"nodeId": gjson.Get(string(msg), "nodeId").Value(),
			}
			p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), resultParams)
		})
	})
	return nil
}

// todo screencast
func (p *protocolAdapter) onStartScreencast(message []byte) []byte {
	params := gjson.Get(string(message), "params")
	format := params.Get("format").String()
	quality := params.Get("quality").Int()
	maxWidth := params.Get("maxWidth").Int()
	maxHeight := params.Get("maxHeight").Int()
	if p.screencast != nil {
		// clear previous session
		p.screencast.stop()
	}
	p.screencast = newScreencastSession(p.adapter,
		WithFormat(format),
		WithMaxWidth(int(maxWidth)),
		WithMaxHeight(int(maxHeight)),
		WithQuality(int(quality)))
	p.screencast.start()

	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), map[string]interface{}{})

	return nil
}

func (p *protocolAdapter) onStopScreencast(message []byte) []byte {
	if p.screencast != nil {
		// clear previous session
		p.screencast.stop()
		p.screencast = nil
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), map[string]interface{}{})

	return nil
}

func (p *protocolAdapter) onScreencastFrameAck(message []byte) []byte {
	if p.screencast != nil {
		frameNumber := gjson.Get(string(message), "params.sessionId").Int()
		// todo Change to int 64?
		p.screencast.ackFrame(int(frameNumber))
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), map[string]interface{}{})

	return nil
}

func (p *protocolAdapter) onGetNavigationHistory(message []byte) []byte {
	var href string
	var id = int(gjson.Get(string(message), "id").Int())
	p.adapter.CallTarget("Runtime.evaluate", map[string]interface{}{"expression": "window.location.href"}, func(result []byte) {
		href = gjson.Get(string(result), "result.value").String()
		p.adapter.CallTarget("Runtime.evaluate", map[string]interface{}{"expression": "window.title"}, func(msg []byte) {
			title := gjson.Get(string(msg), "result.value").String()
			p.adapter.FireResultToTools(id, map[string]interface{}{
				"currentIndex": 0, "entries": []interface{}{
					map[string]interface{}{
						"id":    0,
						"url":   href,
						"title": title,
					},
				},
			})
		})
	})
	return nil
}

func (p *protocolAdapter) onEmulateTouchFromMouseEvent(message []byte) []byte {
	var funcStr = `function simulate(params) {
                const element = document.elementFromPoint(params.x, params.y);
                const e = new MouseEvent(params.type, {
                    screenX: params.x,
                    screenY: params.y,
                    clientX: 0,
                    clientY: 0,
                    ctrlKey: (params.modifiers & 2) === 2,
                    shiftKey: (params.modifiers & 8) === 8,
                    altKey: (params.modifiers & 1) === 1,
                    metaKey: (params.modifiers & 4) === 4,
                    button: params.button,
                    bubbles: true,
                    cancelable: false
                });
                element.dispatchEvent(e);
                return element;
            }`
	newMsg := string(message)
	oldMsg := string(message)
	var err error
	switch gjson.Get(newMsg, "params.type").String() {
	case "mousePressed":
		newMsg, err = sjson.Set(newMsg, "params.type", "mousedown")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		break
	case "mouseReleased":
		newMsg, err = sjson.Set(newMsg, "params.type", "click")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		break
	case "mouseMoved":
		newMsg, err = sjson.Set(newMsg, "params.type", "mousemove")
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		break
	default:
		log.Println(fmt.Sprintf("Unknown emulate mouse event name %s",
			gjson.Get(oldMsg, "params.type")),
		)
	}
	var exp = fmt.Sprintf("(%s)(%s)", funcStr, newMsg)

	p.adapter.CallTarget("Runtime.evaluate", map[string]interface{}{
		"expression": exp,
	}, func(result []byte) {
		if gjson.Get(newMsg, "params.type").String() == "click" {
			newMsg, err = sjson.Set(newMsg, "params.type", "mouseup")
			if err != nil {
				log.Println(e.Convert(err).ToStr())
			}
		}
		p.adapter.CallTarget("Runtime.evaluate", map[string]interface{}{
			"expression": exp,
		}, nil)
	})
	return p.adapter.ReplyWithEmpty(newMsg)
}

func (p *protocolAdapter) onCanEmulateNetworkConditions(message []byte) []byte {
	result := map[string]interface{}{
		"result": false,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *protocolAdapter) onConsoleMessageAdded(message []byte) []byte {
	resultMessage := gjson.Get(string(message), "params.message")
	messageType := resultMessage.Get("type").String()
	var resultType string
	if resultType == "log" {
		switch messageType {
		case "log":
			resultType = "log"
		case "info":
			resultType = "info"
		case "error":
			resultType = "error"
		default:
			resultType = "log"
		}
	} else {
		resultType = "log"
	}
	consoleMessage := map[string]interface{}{
		"source":           gjson.Get(string(message), "source").Value(),
		"level":            resultType,
		"text":             gjson.Get(string(message), "text").Value(),
		"lineNumber":       gjson.Get(string(message), "line").Value(),
		"timestamp":        time.Now().UnixNano(),
		"url":              gjson.Get(string(message), "url").Value(),
		"networkRequestId": gjson.Get(string(message), "networkRequestId").Value(),
	}
	if gjson.Get(string(message), "stackTrace").Exists() {
		consoleMessage["stackTrace"] = map[string]interface{}{
			"callFrames": gjson.Get(string(message), "stackTrace").Value(),
		}
	} else {
		consoleMessage["stackTrace"] = nil
	}
	p.adapter.FireEventToTools("Log.entryAdded", consoleMessage)
	return nil
}

func (p *protocolAdapter) enumerateStyleSheets(message []byte) []byte {
	p.adapter.CallTarget("CSS.getAllStyleSheets", map[string]interface{}{}, func(message []byte) {
		newMsg := string(message)
		oldMsg := string(message)
		var err error
		headers := gjson.Get(string(message), "headers")
		if headers.Exists() {
			for _, header := range headers.Array() {
				newMsg, err = sjson.Set(newMsg, header.Get("isInline").Path(oldMsg), false)
				if err != nil {
					log.Println(e.Convert(err).ToStr())
				}
				newMsg, err = sjson.Set(newMsg, header.Get("startLine").Path(oldMsg), 0)
				if err != nil {
					log.Println(e.Convert(err).ToStr())
				}
				newMsg, err = sjson.Set(newMsg, header.Get("startColumn").Path(oldMsg), 0)
				if err != nil {
					log.Println(e.Convert(err).ToStr())
				}
				p.adapter.FireEventToTools("CSS.styleSheetAdded", map[string]interface{}{
					"header": gjson.Get(newMsg, header.Path(oldMsg)).Value(),
				})
			}
		}
	})
	return nil
}

func (p *protocolAdapter) onAddRule(message []byte) []byte {
	var selector = gjson.Get(string(message), "params.ruleText").String()
	selector = strings.TrimSpace(selector)
	// todo prone to bugs
	selector = strings.Replace(selector, "{}", "", -1)
	params := map[string]interface{}{
		"contextNodeId": p.lastNodeId,
		"selector":      selector,
	}
	p.adapter.CallTarget("CSS.addRule", params, func(addRuleResultMessage []byte) {
		var addRuleResult = &WebKitProtocol.AddRuleResult{}
		err := json.Unmarshal(addRuleResultMessage, addRuleResult)
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}
		p.mapRule(addRuleResult.Rule)

		p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), addRuleResult)
	})
	return nil
}

func (p *protocolAdapter) mapRule(cssRule *WebKitProtocol.CSSRule) {
	if cssRule.RuleId != nil {
		cssRule.DevToolStyleSheetId = cssRule.RuleId.StyleSheetId
		cssRule.RuleId = nil
	}

	p.mapSelectorList(cssRule.SelectorList)

	p.mapStyle(cssRule.Style, cssRule.Origin)

	cssRule.SourceLine = nil
}

func (p *protocolAdapter) onGetMatchedStylesForNodeResult(message []byte) []byte {

	if gjson.Get(string(message), "result").Exists() {
		var getMatchedStylesForNodeResult = &WebKitProtocol.GetMatchedStylesForNodeResult{}

		err := json.Unmarshal([]byte(gjson.Get(string(message), "result").String()), getMatchedStylesForNodeResult)
		if err != nil {
			log.Println(e.Convert(err).ToStr())
		}

		for _, matchedCSSRule := range getMatchedStylesForNodeResult.MatchedCSSRules {
			if matchedCSSRule.Rule != nil {
				p.mapRule(matchedCSSRule.Rule)
			}
		}

		for _, inherited := range getMatchedStylesForNodeResult.Inherited {
			if inherited.MatchedCSSRules != nil {
				for _, matchedCSSRule := range inherited.MatchedCSSRules {
					p.mapRule(matchedCSSRule.Rule)
				}
			}
		}

		newMessage, err1 := sjson.Set(string(message), "result", getMatchedStylesForNodeResult)
		if err1 != nil {
			log.Println(e.Convert(err1).ToStr())
		}
		message = []byte(newMessage)
	}
	return message
}

// onSetStyleTexts todo KeyCheck
func (p *protocolAdapter) onSetStyleTexts(message []byte) []byte {
	var msg = string(message)
	var allStyleText []interface{}
	resultId := gjson.Get(msg, "id").Int()
	editsResult := gjson.Get(msg, "params.edits").Array()
	var whetherToContinueTheCycle = true

	for _, edit := range editsResult {
		if !whetherToContinueTheCycle {
			break
		}
		paramsGetStyleSheet := map[string]interface{}{
			"styleSheetId": edit.Get("styleSheetId").String(),
		}
		p.adapter.CallTarget("CSS.getStyleSheet", paramsGetStyleSheet, func(message []byte) {
			var result = string(message)
			styleSheet := gjson.Get(result, "styleSheet")
			styleSheetRules := gjson.Get(result, "styleSheet.rules")
			if !styleSheet.Exists() || !styleSheetRules.Exists() {
				log.Println("iOS returned a value we were not expecting for getStyleSheet")
			}
			for index, rule := range styleSheetRules.Array() {
				if compareRanges(rule.Get("style.range"), edit.Get("range")) {
					params := map[string]interface{}{
						"styleId": map[interface{}]interface{}{
							"styleSheetId": edit.Get("styleSheetId").String(),
							"ordinal":      index,
						},
						"text": edit.Get("text").String(),
					}
					p.adapter.CallTarget("CSS.setStyleText", params, func(setStyleResult []byte) {
						var setStyleResultData = &WebKitProtocol.SetStyleTextResult{}
						err := json.Unmarshal(setStyleResult, setStyleResultData)
						if err != nil {
							log.Println(e.Convert(err).ToStr())
						}
						p.mapStyle(setStyleResultData.Style, nil)

						allStyleText = append(allStyleText, setStyleResultData.Style)
						// stop for
						whetherToContinueTheCycle = false
					})
				}
			}
		})
	}
	result := map[string]interface{}{
		"styles": allStyleText,
	}
	p.adapter.FireResultToTools(int(resultId), result)
	return nil
}

func (p *protocolAdapter) mapStyle(cssStyle *WebKitProtocol.CSSStyle, ruleOrigin *WebKitProtocol.StyleSheetOrigin) {

	if cssStyle.CssText != nil {
		disabled := p.extractDisabledStyles(*cssStyle.CssText, cssStyle.Range)
		for i, value := range disabled {
			noSpaceStr := strings.TrimSpace(*value.Content)
			// 原版 const text = disabled[i].content.trim().replace(/^\/\*\s*/, '').replace(/;\s*\*\/$/, '');
			reg := regexp.MustCompile(`^\\/\\*\\s*`)
			noSpaceStr = reg.ReplaceAllString(noSpaceStr, ``)

			reg = regexp.MustCompile(`;\\s*\\*\\/$`)
			noSpaceStr = reg.ReplaceAllString(noSpaceStr, ``)

			parts := strings.Split(noSpaceStr, ":")
			if cssStyle.CssProperties != nil {
				var index = len(cssStyle.CssProperties)
				for j, _ := range cssStyle.CssProperties {
					if cssStyle.CssProperties[j].Range != nil &&
						(cssStyle.CssProperties[j].Range.StartLine > disabled[i].Range.StartLine ||
							cssStyle.CssProperties[j].Range.StartLine == disabled[i].Range.StartLine ||
							cssStyle.CssProperties[j].Range.StartColumn > disabled[i].Range.StartColumn) {
						index = j
						break
					}
				}

				cssPropertiesObjects := cssStyle.CssProperties
				// insert index
				cssPropertiesLeft := cssPropertiesObjects[:index+1]
				cssPropertiesRight := cssPropertiesObjects[index+1:]

				implicity := false

				var status WebKitProtocol.CSSPropertyStatus = "disabled"
				data := WebKitProtocol.CSSProperty{
					Implicit: &implicity,
					Name:     &parts[0],
					Range:    disabled[i].Range,
					Status:   &status,
					Text:     disabled[i].Content,
					Value:    &parts[1],
				}

				cssPropertiesLeft = append(cssPropertiesLeft, data)

				cssPropertiesLeft = append(cssPropertiesLeft, cssPropertiesRight...)

				cssStyle.CssProperties = cssPropertiesLeft
			}
		}
	}

	for _, cssProperty := range cssStyle.CssProperties {
		p.mapCssProperty(&cssProperty)
	}
	if *ruleOrigin != "user-agent" {
		cssStyle.StyleSheetId = cssStyle.StyleId.StyleSheetId
		arr, err1 := json.Marshal(cssStyle.Range)
		if err1 != nil {
			log.Println(e.Convert(err1).ToStr())
		}
		var styleKey = fmt.Sprintf("%s_%s", *cssStyle.StyleSheetId, string(arr))
		if p.styleMap == nil {
			p.styleMap = make(map[string]interface{})

		}
		p.styleMap[styleKey] = cssStyle.StyleId

	}
	// delete
	cssStyle.StyleId = nil
	cssStyle.Width = nil
	cssStyle.Height = nil
	// todo         delete cssStyle.sourceLine; this old version?
	//        delete cssStyle.sourceURL;
}

func (p *protocolAdapter) mapCssProperty(cssProperty *WebKitProtocol.CSSProperty) {
	resultTrue := true
	resultFalse := false
	if cssProperty.Status != nil {
		if *cssProperty.Status == "disabled" {
			cssProperty.Disabled = &resultTrue
		} else if *cssProperty.Status == "active" {
			cssProperty.Disabled = &resultFalse
		}
		cssProperty.Status = nil
	}

	priority := cssProperty.Priority
	if priority != nil && *priority != "" {
		cssProperty.Implicit = &resultTrue
	} else {
		cssProperty.Implicit = &resultFalse
	}

	cssProperty.Implicit = nil
}

// extractDisabledStyles todo KeyCheck
func (p *protocolAdapter) extractDisabledStyles(styleText string, cssRange *WebKitProtocol.SourceRange) []WebKitProtocol.SourceRange {
	var startIndices []int
	var styles []WebKitProtocol.SourceRange
	for index, _ := range styleText {
		endIndexBEGINCOMMENT := index + len(BEGIN_COMMENT)
		endIndexENDCOMMENT := index + len(END_COMMENT)
		if endIndexBEGINCOMMENT <= len(styleText) && string([]rune(styleText)[index:endIndexBEGINCOMMENT]) == BEGIN_COMMENT {
			startIndices = append(startIndices, index)
			index = index + len(BEGIN_COMMENT)
		} else if endIndexENDCOMMENT <= len(styleText) && string([]rune(styleText)[index:endIndexENDCOMMENT]) == END_COMMENT {
			if len(startIndices) == 0 {
				return []WebKitProtocol.SourceRange{}
			}
			startIndex := startIndices[0]
			startIndices = startIndices[1:]
			endIndex := index + len(END_COMMENT)

			startRangeLine, startRangeColumn := p.getLineColumnFromIndex(styleText, startIndex, cssRange)
			endRangeLine, endRangeColumn := p.getLineColumnFromIndex(styleText, endIndex, cssRange)

			content := styleText[startIndex:endIndex]
			propertyItem := WebKitProtocol.SourceRange{
				Content: &content,
				Range: &WebKitProtocol.SourceRange{
					StartLine:   startRangeLine,
					StartColumn: startRangeColumn,
					EndLine:     endRangeLine,
					EndColumn:   endRangeColumn,
				},
			}
			styles = append(styles, propertyItem)
			index = endIndex - 1
		}
	}
	if len(startIndices) == 0 {
		return []WebKitProtocol.SourceRange{}
	}
	return styles
}

// todo KeyCheck
func (p *protocolAdapter) getLineColumnFromIndex(text string, index int, startRange *WebKitProtocol.SourceRange) (line int, column int) {
	if text == "" || index < 0 || index > len(text) {
		return 0, 0
	}
	if startRange != nil {
		line = startRange.StartLine
		column = startRange.StartColumn
	}
	for i := 0; i < len(text) && i < index; i++ {
		if text[i] == '\r' && i+1 < len(text) && text[i+1] == '\n' {
			i++
			line++
			column = 0
		} else if text[i] == '\n' || text[i] == '\r' {
			line++
			column = 0
		} else {
			column++
		}
	}
	return line, column
}

func compareRanges(rangeLeft gjson.Result, rangeRight gjson.Result) bool {
	return rangeLeft.Get("startLine").Int() == rangeRight.Get("startLine").Int() &&
		rangeLeft.Get("startColumn").Int() == rangeRight.Get("startColumn").Int() &&
		rangeLeft.Get("endLine").Int() == rangeRight.Get("endLine").Int() &&
		rangeLeft.Get("endColumn").Int() == rangeRight.Get("endColumn").Int()
}

func ReplaceMethodNameAndOutputBinary(message []byte, method string) []byte {
	var msg = make(map[string]interface{})
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	// todo Regular?
	msg["method"] = method

	arr, err1 := json.Marshal(msg)
	if err1 != nil {
		log.Println(e.Convert(err1).ToStr())
	}
	return arr
}

var BEGIN_COMMENT = "/* "
var END_COMMENT = " */"
var SEPARATOR = ": "

/******/ (() => { // webpackBootstrap
/******/ 	var __webpack_modules__ = ({

/***/ "./src/app.ts":
/*!********************!*\
  !*** ./src/app.ts ***!
  \********************/
/***/ (() => {

function _typeof(o) { "@babel/helpers - typeof"; return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (o) { return typeof o; } : function (o) { return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o; }, _typeof(o); }
function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _unsupportedIterableToArray(arr) || _nonIterableSpread(); }
function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }
function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }
function _iterableToArray(iter) { if (typeof Symbol !== "undefined" && iter[Symbol.iterator] != null || iter["@@iterator"] != null) return Array.from(iter); }
function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) return _arrayLikeToArray(arr); }
function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) arr2[i] = arr[i]; return arr2; }
function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }
function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, _toPropertyKey(descriptor.key), descriptor); } }
function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); Object.defineProperty(Constructor, "prototype", { writable: false }); return Constructor; }
function _defineProperty(obj, key, value) { key = _toPropertyKey(key); if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }
function _toPropertyKey(t) { var i = _toPrimitive(t, "string"); return "symbol" == _typeof(i) ? i : String(i); }
function _toPrimitive(t, r) { if ("object" != _typeof(t) || !t) return t; var e = t[Symbol.toPrimitive]; if (void 0 !== e) { var i = e.call(t, r || "default"); if ("object" != _typeof(i)) return i; throw new TypeError("@@toPrimitive must return a primitive value."); } return ("string" === r ? String : Number)(t); }
var ActionMode;
(function (ActionMode) {
  ActionMode[ActionMode["Pencil"] = 0] = "Pencil";
  ActionMode[ActionMode["Select"] = 1] = "Select";
  ActionMode[ActionMode["Move"] = 2] = "Move";
})(ActionMode || (ActionMode = {}));
var PixelEditor = /*#__PURE__*/function () {
  function PixelEditor(canvas) {
    _classCallCheck(this, PixelEditor);
    _defineProperty(this, "canvas", void 0);
    _defineProperty(this, "ctx", void 0);
    _defineProperty(this, "originalWidth", 64);
    _defineProperty(this, "originalHeight", 32);
    _defineProperty(this, "padWidth", 32);
    _defineProperty(this, "padHeight", 16);
    _defineProperty(this, "canvasScale", 10);
    _defineProperty(this, "loadedSpritesheet", void 0);
    _defineProperty(this, "colorsDiv", void 0);
    _defineProperty(this, "framesDiv", void 0);
    _defineProperty(this, "selectedColor", void 0);
    _defineProperty(this, "selectedFrame", void 0);
    _defineProperty(this, "actionMode", void 0);
    this.canvas = canvas;
    this.ctx = this.canvas.getContext('2d');
  }
  _createClass(PixelEditor, [{
    key: "addControls",
    value: function addControls() {
      var _this = this;
      var functions = {
        'savespritesheet': function savespritesheet(e) {},
        // this.saveSpritesheet(),
        'move-first': function moveFirst(e) {},
        // this.moveFrameFirst(),
        'move-previous': function movePrevious(e) {},
        // this.moveFramePrevious(),
        'move-next': function moveNext(e) {},
        // this.moveFrameNext(),
        'move-last': function moveLast(e) {} // this.moveFrameLast(),
        // 'play-animation': (e: Event) => { }, // this.playAnimation(),
      };
      Object.keys(functions).map(function (id) {
        var el = document.getElementById(id);
        el.addEventListener('click', function (e) {
          e.preventDefault();
          functions[id](e);
        });
      });
      var loadSpritesheetBtn = document.getElementById('loadspritesheet');
      var loadSpritesheetInput = document.getElementById('loadspritesheetinput');
      loadSpritesheetBtn.addEventListener('click', function () {
        return loadSpritesheetInput.click();
      });
      loadSpritesheetInput.addEventListener('change', function (e) {
        var file = e.target.files[0];
        if (!file) return;
        var reader = new FileReader();
        reader.onload = function (event) {
          _this.loadSpritesheet(event.target.result);
        };
        reader.readAsText(file);
      });
    }
  }, {
    key: "loadSpritesheet",
    value: function loadSpritesheet(filecontents) {
      var spritesheetJson;
      try {
        spritesheetJson = JSON.parse(filecontents);
      } catch (err) {
        console.error(err);
        return;
      }
      this.loadedSpritesheet = {
        colors: spritesheetJson.colors || [0],
        pixeldata: spritesheetJson.pixeldata || [[]],
        animation: spritesheetJson.animation || [],
        width: spritesheetJson.width || 0,
        height: spritesheetJson.height || 0,
        num_sheets: spritesheetJson.num_sheets || 0,
        fps: spritesheetJson.fps || 0
      };
      this.createColors();
      this.createFrames();
    }
  }, {
    key: "colorToHex",
    value: function colorToHex(color) {
      return "#".concat(color.toString(16).padStart(6, '0'));
    }
  }, {
    key: "createColors",
    value: function createColors() {
      var _this2 = this;
      if (!this.colorsDiv) this.colorsDiv = document.getElementById('colors');
      _toConsumableArray(this.colorsDiv.querySelectorAll('div')).forEach(function (child) {
        return child.parentElement.removeChild(child);
      });
      this.loadedSpritesheet.colors.forEach(function (color, colorIndex) {
        var colorDiv = document.createElement('div');
        colorDiv.classList.add('color');
        colorDiv.setAttribute('data-color', color.toString());
        colorDiv.addEventListener('click', function (e) {
          _this2.selectedColor = parseInt(colorDiv.getAttribute('data-color'), 10);
          _toConsumableArray(_this2.colorsDiv.querySelectorAll('.selected')).forEach(function (el) {
            return el.classList.remove('selected');
          });
          colorDiv.classList.add('selected');
        });
        var hexColor = _this2.colorToHex(color);
        var colorSample = document.createElement('div');
        colorSample.classList.add('color-sample');
        colorSample.style.background = hexColor;
        colorDiv.appendChild(colorSample);
        var colorInput = document.createElement('input');
        colorInput.type = 'text';
        colorInput.classList.add('color-input');
        colorInput.value = hexColor;
        colorInput.addEventListener('blur', function (e) {
          var target = e.target;
          if (target.value === hexColor) return;
          var hexRegex = /^#[0-9a-fA-f]{3,6}$/;
          if (hexRegex.test(target.value) !== true) {
            target.value = hexColor;
            return;
          }
          _this2.loadedSpritesheet.colors[colorIndex] = parseInt(target.value.slice(1), 16);
          _this2.createColors();
        });
        colorDiv.appendChild(colorInput);
        _this2.colorsDiv.appendChild(colorDiv);
      });
      this.colorsDiv.querySelectorAll('div')[0].classList.add('selected');
      this.selectedColor = this.loadedSpritesheet.colors[0];
    }
  }, {
    key: "createFrames",
    value: function createFrames() {
      var _this3 = this;
      if (!this.framesDiv) this.framesDiv = document.getElementById('frames');
      _toConsumableArray(this.framesDiv.querySelectorAll('div')).forEach(function (child) {
        return child.parentElement.removeChild(child);
      });
      this.loadedSpritesheet.pixeldata.forEach(function (sheet, sheetIndex) {
        var frameDiv = document.createElement('div');
        frameDiv.classList.add('frame');
        frameDiv.innerHTML = (sheetIndex + 1).toString();
        frameDiv.addEventListener('click', function (e) {
          _this3.selectedFrame = sheetIndex;
          _toConsumableArray(_this3.framesDiv.querySelectorAll('.selected')).forEach(function (el) {
            return el.classList.remove('selected');
          });
          frameDiv.classList.add('selected');
          _this3.drawFrame();
        });
        _this3.framesDiv.appendChild(frameDiv);
      });
      this.framesDiv.querySelectorAll('div')[0].classList.add('selected');
      this.selectedFrame = 0;
      this.drawFrame();
    }
  }, {
    key: "resizeCanvas",
    value: function resizeCanvas() {
      this.canvas.width = this.originalWidth * this.canvasScale + this.padWidth * this.canvasScale;
      this.canvas.height = this.originalHeight * this.canvasScale + this.padHeight * this.canvasScale;
    }
  }, {
    key: "clearCanvas",
    value: function clearCanvas() {
      var originalFillStyle = this.ctx.fillStyle;
      this.ctx.fillStyle = '#fff';
      this.ctx.beginPath();
      this.ctx.rect(0, 0, this.canvas.width, this.canvas.height);
      this.ctx.fill();
      this.ctx.fillStyle = originalFillStyle;
    }
  }, {
    key: "drawWindow",
    value: function drawWindow() {
      var originalStrokeStyle = this.ctx.strokeStyle;
      var x1 = 0 + this.padWidth * this.canvasScale / 2;
      var y1 = 0 + this.padHeight * this.canvasScale / 2;
      var x2 = this.originalWidth * this.canvasScale;
      var y2 = this.originalHeight * this.canvasScale;
      this.ctx.beginPath();
      this.ctx.rect(x1, y1, x2, y2);
      this.ctx.lineWidth = 1;
      this.ctx.stroke();
      this.ctx.strokeStyle = originalStrokeStyle;
    }
  }, {
    key: "drawGrid",
    value: function drawGrid() {
      var originalStrokeStyle = this.ctx.strokeStyle;
      this.ctx.strokeStyle = '#ccc';
      this.ctx.beginPath();
      for (var i = 0; i <= this.canvas.height; i += this.canvasScale) {
        this.ctx.moveTo(0, i);
        this.ctx.lineTo(this.canvas.width, i);
      }
      for (var j = 0; j <= this.canvas.width; j += this.canvasScale) {
        this.ctx.moveTo(j, 0);
        this.ctx.lineTo(j, this.canvas.height);
      }
      this.ctx.stroke();
      this.ctx.strokeStyle = originalStrokeStyle;
    }
  }, {
    key: "addCanvasListeners",
    value: function addCanvasListeners() {
      var _this4 = this;
      this.canvas.addEventListener('click', function (e) {
        var x = Math.floor((e.clientX - _this4.canvas.getBoundingClientRect().left) / _this4.canvasScale) - _this4.padWidth / 2;
        var y = Math.floor((e.clientY - _this4.canvas.getBoundingClientRect().top) / _this4.canvasScale) - _this4.padHeight / 2;
        _this4.setPixel(x, y, _this4.selectedColor);
      });
    }
  }, {
    key: "drawFrame",
    value: function drawFrame() {
      this.clearCanvas();
      this.drawGrid();
      this.drawWindow();
      var sheet = this.loadedSpritesheet.pixeldata[this.selectedFrame];
      for (var y = 0; y < sheet.length; y++) {
        for (var x = 0; x < sheet[y].length; x++) {
          var color = this.loadedSpritesheet.colors[sheet[y][x]];
          if (color != undefined) {
            this.setPixel(x, y, color);
          }
        }
      }
    }
  }, {
    key: "setPixel",
    value: function setPixel(x, y, color) {
      var hexColor = this.colorToHex(color);
      var realX = (x + this.padWidth / 2) * this.canvasScale;
      var realY = (y + this.padHeight / 2) * this.canvasScale;
      var originalFillStyle = this.ctx.fillStyle;
      this.ctx.fillStyle = hexColor;
      this.ctx.beginPath();
      this.ctx.rect(realX, realY, this.canvasScale, this.canvasScale);
      this.ctx.fill();
      this.ctx.fillStyle = originalFillStyle;
    }
  }]);
  return PixelEditor;
}();
window.addEventListener('DOMContentLoaded', function () {
  var canvas = document.getElementById('canvas');
  var pe = new PixelEditor(canvas);
  pe.resizeCanvas();
  pe.addControls();
  pe.addCanvasListeners();
  pe.loadSpritesheet('{}');
});

/***/ }),

/***/ "./src/styles.scss":
/*!*************************!*\
  !*** ./src/styles.scss ***!
  \*************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

"use strict";
__webpack_require__.r(__webpack_exports__);
// extracted by mini-css-extract-plugin


/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = __webpack_modules__;
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/chunk loaded */
/******/ 	(() => {
/******/ 		var deferred = [];
/******/ 		__webpack_require__.O = (result, chunkIds, fn, priority) => {
/******/ 			if(chunkIds) {
/******/ 				priority = priority || 0;
/******/ 				for(var i = deferred.length; i > 0 && deferred[i - 1][2] > priority; i--) deferred[i] = deferred[i - 1];
/******/ 				deferred[i] = [chunkIds, fn, priority];
/******/ 				return;
/******/ 			}
/******/ 			var notFulfilled = Infinity;
/******/ 			for (var i = 0; i < deferred.length; i++) {
/******/ 				var [chunkIds, fn, priority] = deferred[i];
/******/ 				var fulfilled = true;
/******/ 				for (var j = 0; j < chunkIds.length; j++) {
/******/ 					if ((priority & 1 === 0 || notFulfilled >= priority) && Object.keys(__webpack_require__.O).every((key) => (__webpack_require__.O[key](chunkIds[j])))) {
/******/ 						chunkIds.splice(j--, 1);
/******/ 					} else {
/******/ 						fulfilled = false;
/******/ 						if(priority < notFulfilled) notFulfilled = priority;
/******/ 					}
/******/ 				}
/******/ 				if(fulfilled) {
/******/ 					deferred.splice(i--, 1)
/******/ 					var r = fn();
/******/ 					if (r !== undefined) result = r;
/******/ 				}
/******/ 			}
/******/ 			return result;
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/jsonp chunk loading */
/******/ 	(() => {
/******/ 		// no baseURI
/******/ 		
/******/ 		// object to store loaded and loading chunks
/******/ 		// undefined = chunk not loaded, null = chunk preloaded/prefetched
/******/ 		// [resolve, reject, Promise] = chunk loading, 0 = chunk loaded
/******/ 		var installedChunks = {
/******/ 			"/public/app": 0,
/******/ 			"public/styles": 0
/******/ 		};
/******/ 		
/******/ 		// no chunk on demand loading
/******/ 		
/******/ 		// no prefetching
/******/ 		
/******/ 		// no preloaded
/******/ 		
/******/ 		// no HMR
/******/ 		
/******/ 		// no HMR manifest
/******/ 		
/******/ 		__webpack_require__.O.j = (chunkId) => (installedChunks[chunkId] === 0);
/******/ 		
/******/ 		// install a JSONP callback for chunk loading
/******/ 		var webpackJsonpCallback = (parentChunkLoadingFunction, data) => {
/******/ 			var [chunkIds, moreModules, runtime] = data;
/******/ 			// add "moreModules" to the modules object,
/******/ 			// then flag all "chunkIds" as loaded and fire callback
/******/ 			var moduleId, chunkId, i = 0;
/******/ 			if(chunkIds.some((id) => (installedChunks[id] !== 0))) {
/******/ 				for(moduleId in moreModules) {
/******/ 					if(__webpack_require__.o(moreModules, moduleId)) {
/******/ 						__webpack_require__.m[moduleId] = moreModules[moduleId];
/******/ 					}
/******/ 				}
/******/ 				if(runtime) var result = runtime(__webpack_require__);
/******/ 			}
/******/ 			if(parentChunkLoadingFunction) parentChunkLoadingFunction(data);
/******/ 			for(;i < chunkIds.length; i++) {
/******/ 				chunkId = chunkIds[i];
/******/ 				if(__webpack_require__.o(installedChunks, chunkId) && installedChunks[chunkId]) {
/******/ 					installedChunks[chunkId][0]();
/******/ 				}
/******/ 				installedChunks[chunkId] = 0;
/******/ 			}
/******/ 			return __webpack_require__.O(result);
/******/ 		}
/******/ 		
/******/ 		var chunkLoadingGlobal = self["webpackChunk"] = self["webpackChunk"] || [];
/******/ 		chunkLoadingGlobal.forEach(webpackJsonpCallback.bind(null, 0));
/******/ 		chunkLoadingGlobal.push = webpackJsonpCallback.bind(null, chunkLoadingGlobal.push.bind(chunkLoadingGlobal));
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module depends on other loaded chunks and execution need to be delayed
/******/ 	__webpack_require__.O(undefined, ["public/styles"], () => (__webpack_require__("./src/app.ts")))
/******/ 	var __webpack_exports__ = __webpack_require__.O(undefined, ["public/styles"], () => (__webpack_require__("./src/styles.scss")))
/******/ 	__webpack_exports__ = __webpack_require__.O(__webpack_exports__);
/******/ 	
/******/ })()
;
interface Spritesheet {
    colors: number[],
    pixeldata: number[][][],
    animation: number[],
    width: number,
    height: number,
    num_sheets: number,
    fps: number,
}

enum ActionMode {
    Pencil,
    Select,
    Move,
}

class PixelEditor {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;

    originalWidth: number = 64;
    originalHeight: number = 32;

    padWidth: number = 32;
    padHeight: number = 16;

    canvasScale: number = 10;

    loadedSpritesheet: Spritesheet;

    colorsDiv: HTMLDivElement;
    framesDiv: HTMLDivElement;

    selectedColor: number;
    selectedFrame: number;

    actionMode: ActionMode;

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
        this.ctx = this.canvas.getContext('2d');
    }

    addControls() {
        const functions = {
            'savespritesheet': (e: Event) => { }, // this.saveSpritesheet(),
            'move-first': (e: Event) => { }, // this.moveFrameFirst(),
            'move-previous': (e: Event) => { }, // this.moveFramePrevious(),
            'move-next': (e: Event) => { }, // this.moveFrameNext(),
            'move-last': (e: Event) => { }, // this.moveFrameLast(),
            // 'play-animation': (e: Event) => { }, // this.playAnimation(),
        }

        Object.keys(functions).map(id => {
            const el = document.getElementById(id);
            el.addEventListener('click', e => {
                e.preventDefault();

                functions[id](e);
            });
        });

        const loadSpritesheetBtn = document.getElementById('loadspritesheet');
        const loadSpritesheetInput = document.getElementById('loadspritesheetinput');
        loadSpritesheetBtn.addEventListener('click', () => loadSpritesheetInput.click());

        loadSpritesheetInput.addEventListener('change', e => {
            const file = (<HTMLInputElement>e.target).files[0];

            if (!file) return;

            const reader = new FileReader();
            reader.onload = event => {
                this.loadSpritesheet(<string>event.target.result);
            }

            reader.readAsText(file);
        });
    }

    loadSpritesheet(filecontents: string) {
        let spritesheetJson: Record<string, any>;

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
            fps: spritesheetJson.fps || 0,
        }

        this.createColors();
        this.createFrames();
    }

    colorToHex(color: number): string {
        return `#${color.toString(16).padStart(6, '0')}`;
    }

    createColors() {
        if (!this.colorsDiv) this.colorsDiv = <HTMLDivElement>document.getElementById('colors');
        [...this.colorsDiv.querySelectorAll('div')].forEach(child => child.parentElement.removeChild(child));

        this.loadedSpritesheet.colors.forEach((color, colorIndex) => {
            const colorDiv = document.createElement('div');

            colorDiv.classList.add('color');
            colorDiv.setAttribute('data-color', color.toString());

            colorDiv.addEventListener('click', e => {
                this.selectedColor = parseInt(colorDiv.getAttribute('data-color'), 10);

                [...this.colorsDiv.querySelectorAll('.selected')].forEach(el => el.classList.remove('selected'));
                colorDiv.classList.add('selected');
            });

            const hexColor = this.colorToHex(color);
            const colorSample = document.createElement('div');

            colorSample.classList.add('color-sample');
            colorSample.style.background = hexColor;
            colorDiv.appendChild(colorSample);

            const colorInput = document.createElement('input');
            colorInput.type = 'text';
            colorInput.classList.add('color-input');
            colorInput.value = hexColor;

            colorInput.addEventListener('blur', e => {
                const target = (<HTMLInputElement>e.target);

                if (target.value === hexColor) return;

                const hexRegex = /^#[0-9a-fA-f]{3,6}$/;
                if (hexRegex.test(target.value) !== true) {
                    target.value = hexColor;
                    return;
                }

                this.loadedSpritesheet.colors[colorIndex] = parseInt(target.value.slice(1), 16);

                this.createColors();
            });

            colorDiv.appendChild(colorInput);

            this.colorsDiv.appendChild(colorDiv);
        });

        this.colorsDiv.querySelectorAll('div')[0].classList.add('selected');
        this.selectedColor = this.loadedSpritesheet.colors[0];
    }

    createFrames() {
        if (!this.framesDiv) this.framesDiv = <HTMLDivElement>document.getElementById('frames');
        [...this.framesDiv.querySelectorAll('div')].forEach(child => child.parentElement.removeChild(child));

        this.loadedSpritesheet.pixeldata.forEach((sheet, sheetIndex) => {
            const frameDiv = document.createElement('div');
            frameDiv.classList.add('frame');
            frameDiv.innerHTML = (sheetIndex + 1).toString();

            frameDiv.addEventListener('click', e => {
                this.selectedFrame = sheetIndex;

                [...this.framesDiv.querySelectorAll('.selected')].forEach(el => el.classList.remove('selected'));
                frameDiv.classList.add('selected');

                this.drawFrame();
            });

            this.framesDiv.appendChild(frameDiv);
        });

        this.framesDiv.querySelectorAll('div')[0].classList.add('selected');
        this.selectedFrame = 0;

        this.drawFrame();
    }

    resizeCanvas() {
        this.canvas.width = (this.originalWidth * this.canvasScale) + (this.padWidth * this.canvasScale);
        this.canvas.height = (this.originalHeight * this.canvasScale) + (this.padHeight * this.canvasScale);
    }

    clearCanvas() {
        const originalFillStyle = this.ctx.fillStyle;

        this.ctx.fillStyle = '#fff';

        this.ctx.beginPath();
        this.ctx.rect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.fill();

        this.ctx.fillStyle = originalFillStyle;
    }

    drawWindow() {
        const originalStrokeStyle = this.ctx.strokeStyle;

        const x1 = 0 + (this.padWidth * this.canvasScale) / 2;
        const y1 = 0 + (this.padHeight * this.canvasScale) / 2;
        const x2 = this.originalWidth * this.canvasScale;
        const y2 = this.originalHeight * this.canvasScale;

        this.ctx.beginPath();
        this.ctx.rect(x1, y1, x2, y2);
        this.ctx.lineWidth = 1;
        this.ctx.stroke();

        this.ctx.strokeStyle = originalStrokeStyle;
    }

    drawGrid() {
        const originalStrokeStyle = this.ctx.strokeStyle;

        this.ctx.strokeStyle = '#ccc';

        this.ctx.beginPath();
        for (let i = 0; i <= this.canvas.height; i += this.canvasScale) {
            this.ctx.moveTo(0, i);
            this.ctx.lineTo(this.canvas.width, i);
        }

        for (let j = 0; j <= this.canvas.width; j += this.canvasScale) {
            this.ctx.moveTo(j, 0);
            this.ctx.lineTo(j, this.canvas.height);
        }
        this.ctx.stroke();

        this.ctx.strokeStyle = originalStrokeStyle;
    }

    addCanvasListeners() {
        this.canvas.addEventListener('click', e => {
            const x = Math.floor((e.clientX - this.canvas.getBoundingClientRect().left) / this.canvasScale) - (this.padWidth / 2);
            const y = Math.floor((e.clientY - this.canvas.getBoundingClientRect().top) / this.canvasScale) - (this.padHeight / 2);

            this.setPixel(x, y, this.selectedColor);
        });
    }

    drawFrame() {
        this.clearCanvas();
        this.drawGrid();
        this.drawWindow();

        const sheet = this.loadedSpritesheet.pixeldata[this.selectedFrame];

        for (let y = 0; y < sheet.length; y++) {
            for (let x = 0; x < sheet[y].length; x++) {
                const color = this.loadedSpritesheet.colors[sheet[y][x]];

                if (color != undefined) {
                    this.setPixel(x, y, color);
                }
            }
        }
    }

    setPixel(x: number, y: number, color: number) {
        const hexColor = this.colorToHex(color);

        const realX = (x + (this.padWidth / 2)) * this.canvasScale;
        const realY = (y + (this.padHeight / 2)) * this.canvasScale;

        const originalFillStyle = this.ctx.fillStyle;

        this.ctx.fillStyle = hexColor;

        this.ctx.beginPath();
        this.ctx.rect(realX, realY, this.canvasScale, this.canvasScale);
        this.ctx.fill();

        this.ctx.fillStyle = originalFillStyle;
    }
}

window.addEventListener('DOMContentLoaded', () => {
    const canvas = <HTMLCanvasElement>document.getElementById('canvas');
    const pe = new PixelEditor(canvas);

    pe.resizeCanvas();

    pe.addControls();
    pe.addCanvasListeners();

    pe.loadSpritesheet('{}');
});

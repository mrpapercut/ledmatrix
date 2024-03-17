const fs = require('fs');
const path = require('path');

const { PNG } = require('pngjs');

const dirname = './mario-spritesheets/individual/mario-smw-walk/';

const results = {
    width: 0,
    height: 0,
    num_sheets: 0,
    fps: 12,
    animation: [],
    colors: [],
    pixeldata: []
}

fs.readdir(dirname, (err, files) => {
    if (err) return;

    files.forEach(file => {
        const filepath = path.join(dirname, file);

        const filecontents = fs.readFileSync(filepath);

        const png = PNG.sync.read(filecontents);

        const sheetPixels = [];

        if (png.width > results.width) results.width = png.width;
        if (png.height > results.height) results.height = png.height;
        results.num_sheets++;

        for (let y = 0; y < png.height; y++) {
            sheetPixels[y] = [];

            for (let x = 0; x < png.width; x++) {
                const index = (png.width * y + x) << 2;
                const [red, green, blue, alpha] = new Array(4).fill(0).map((_, i) => png.data[index + i]);

                if (alpha === 0) {
                    sheetPixels[y][x] = -1;
                } else {
                    const hexColor = eval(`0x${[red, green, blue].map(c => c.toString(16).padStart(2, '0')).join('')}`);

                    if (!results.colors.includes(hexColor)) results.colors.push(hexColor);

                    sheetPixels[y][x] = results.colors.indexOf(hexColor);
                }
            }
        }

        results.pixeldata.push(sheetPixels);
    });

    results.animation = new Array(results.pixeldata.length).fill().map((_, i) => i);

    console.log(JSON.stringify(results));
});

function convertObjToCPP(spritesheet) {
    const { width, height, numSheets, colors, pixelData } = spritesheet;
    const sheets = [];

    pixelData.forEach((sheet, idx) => {
        sheets[idx] = `{${sheet.map(row => `{${row.join(', ')}}`).join(',\n')}}`;
    });

    return `{\n${width}, ${height}, ${numSheets},\n{${colors.map(c => `0x${c.toString(16).padStart(6, '0')}`).join(', ')}},\n{${sheets.join(',\n')}}}`;
}

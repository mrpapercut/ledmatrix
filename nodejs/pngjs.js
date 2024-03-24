const fs = require('fs/promises');
const path = require('path');

const { PNG } = require('pngjs');

// const locationname = './mario-spritesheets/individual/mario-smw-walk/';
const locationname = './mario-spritesheets/individual/yt-logos/grandpoobear.png';

const results = {
    width: 0,
    height: 0,
    num_sheets: 0,
    fps: 12,
    animation: [],
    colors: [0],
    pixeldata: []
}

async function processFiles() {
    const stat = await fs.lstat(locationname)
    const isDir = stat.isDirectory();

    if (isDir) {
        try {
            const dirfiles = await fs.readdir(locationname);

            for (let i = 0; i < dirfiles.length; i++) {
                const filepath = path.join(locationname, dirfiles[i]);

                const filecontents = await fs.readFile(filepath);

                handleFile(filecontents)
            }
        } catch (err) {
            return;
        }
    } else {
        results.fps = 1;

        const filecontents = await fs.readFile(locationname);

        handleFile(filecontents);
    }

    results.animation = new Array(results.pixeldata.length).fill().map((_, i) => i);

    console.log(JSON.stringify(results));
}

(async () => {
    processFiles()
})();

function handleFile(filecontents) {
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
}

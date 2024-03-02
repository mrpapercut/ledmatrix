#pragma once

#include <vector>
#include <string>

typedef std::vector<std::vector<std::vector<int>>> PixelData;
typedef std::vector<std::vector<int>> FontData;

class Spritesheet
{
public:
    int width;
    int height;
    int numSheets;
    std::vector<int> colors;
    PixelData pixelData;

    Spritesheet(int width, int height, int numSheets, std::vector<int> colors, PixelData pixelData);
    virtual ~Spritesheet() {}

    static Spritesheet ConvertFontToPixels(FontData fontData);
    static Spritesheet ConvertCharactersToSpritesheet(FontData fontData);
    static Spritesheet CombineSpritesheets(std::vector<Spritesheet> spritesheets);
};

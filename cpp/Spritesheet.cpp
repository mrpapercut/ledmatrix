#include <bitset>
#include "Spritesheet.hpp"

Spritesheet::Spritesheet(int width, int height, int numSheets = 0, std::vector<int> colors = {}, PixelData pixelData = {})
    : width(width), height(height), numSheets(numSheets), colors(std::move(colors)), pixelData(std::move(pixelData)) {}

Spritesheet Spritesheet::ConvertFontToPixels(FontData fontData)
{
    Spritesheet textPixels = Spritesheet(0, 0, 1, {0x0, 0x44aa00}, {});

    // int maxWidth = 0;

    for (size_t i = 0; i < fontData.size(); ++i) // Per letter
    {
        std::vector<std::vector<int>> characterPixels = {};

        for (size_t j = 0; j < fontData[i].size(); ++j) // Per row in letter
        {
            std::string binaryString = std::bitset<sizeof(int) * 8>(fontData[i][j]).to_string();
            std::string reversedString = std::string(binaryString.rbegin(), binaryString.rend());

            reversedString.erase(reversedString.find_last_of('1') + 1);

            std::vector<int> booleanList;
            for (char bit : reversedString) // Per pixel
            {
                booleanList.push_back(bit == '1' ? 1 : 0);
            }

            characterPixels.push_back(booleanList);
        }

        textPixels.pixelData.push_back(characterPixels);
    }

    return textPixels;
};

Spritesheet Spritesheet::ConvertCharactersToSpritesheet(FontData fontData)
{
    Spritesheet emptySpritesheet = Spritesheet(0, 0, 1, {}, {});

    return emptySpritesheet;
}

Spritesheet Spritesheet::CombineSpritesheets(std::vector<Spritesheet> spritesheets)
{
    Spritesheet emptySpritesheet = Spritesheet(0, 0, 1, {}, {});

    return emptySpritesheet;
}

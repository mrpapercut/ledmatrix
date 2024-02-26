#include "led-matrix.h"

#include <iostream>
#include <fstream>
#include <vector>
#include <bitset>

#include <unistd.h>
#include <math.h>
#include <stdio.h>
#include <signal.h>

#include "constants.h"

#include "fonts.hpp"
#include "Text.hpp"
#include "spritesheets.hpp"

namespace Text
{
    Text::Text(Canvas *canvas, volatile bool &interrupt_flag) : canvas_(canvas), interrupt_received(interrupt_flag) {}

    void Text::ScrollText(Font font, const char *message)
    {
        int fps = 12;

        int offsetX = -24;

        while (!interrupt_received)
        {
            canvas()->Clear();

            offsetX = offsetX + 3;
            if (offsetX > (SCREEN_WIDTH + 24))
            {
                offsetX = 0 - 24;
            }

            usleep((1000 / fps) * 1000);
        }
    }

    void Text::ShowText(Font font, const char *message)
    {
        std::vector<std::vector<int>> characterData = {};

        for (int i = 0; message[i] != '\0'; ++i)
        {
            char currentChar = message[i];

            if (font.characters.count(currentChar) > 0)
            {
                characterData.push_back(font.characters[currentChar]);
            }
        }

        Spritesheet spritesheet = ConvertFontToPixels(characterData);

        std::vector<int> colors = spritesheet.colors;

        int fps = 1;

        while (!interrupt_received)
        {
            int offsetX = 0; // 0 - maxSpriteWidth;
            int offsetY = 0; // (SCREEN_HEIGHT - maxSpriteHeight) / 2;

            size_t characterWidth;

            canvas()->Clear();

            // For text, a single sheet within pixelData is a single character
            for (size_t charIndex = 0; charIndex < spritesheet.pixelData.size(); ++charIndex) // Character
            {
                std::vector<std::vector<int>> currentCharacter = spritesheet.pixelData[charIndex];
                characterWidth = 0;

                for (size_t y = 0; y < currentCharacter.size(); ++y) // Row
                {
                    for (size_t x = 0; x < currentCharacter[y].size(); ++x) // Pixel
                    {

                        int colorIndex = currentCharacter[y][x];
                        if (colorIndex == 0)
                            continue;

                        int color = colors[colorIndex];
                        if (color == 0)
                            continue;

                        if (x > characterWidth)
                        {
                            characterWidth = x;
                        }

                        int red = (color >> 16) & 0xff;
                        int green = (color >> 8) & 0xff;
                        int blue = color & 0xff;

                        canvas()->SetPixel(x + offsetX, y + offsetY, red, green, blue);
                    }
                }

                offsetX += characterWidth + 2;

                if (offsetX + characterWidth > SCREEN_WIDTH)
                {
                    offsetX = 0;
                    offsetY += 10;
                }
            }

            usleep((1000 / fps) * 1000);
        }
    }

    Spritesheet Text::ConvertFontToPixels(std::vector<std::vector<int>> fontData)
    {
        Spritesheet textPixels = {0, 0, 1, {0x0, 0x44aa00}, {}};

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
    }

    class ShowHelloWorld : public Text
    {
    public:
        Font font = Fonts::DefaultFont;

        ShowHelloWorld(Canvas *canvas, volatile bool &interrupt_flag) : Text(canvas, interrupt_flag) {}

        void Run() override
        {
            const char *message = "Hello, world!";
            ShowText(font, message);
        }
    };

    class ShowRandomLineFromFile : public Text
    {
    public:
        Font font = Fonts::DefaultFont;

        ShowRandomLineFromFile(Canvas *canvas, volatile bool &interrupt_flag) : Text(canvas, interrupt_flag) {}

        void Run() override
        {
            std::string message = ReadLineFromFile();

            ShowText(font, message.c_str());
        }

    private:
        std::string ReadLineFromFile()
        {
            std::ifstream file("inputdata.txt");

            if (!file.is_open())
            {
                std::perror("Unable to open inputdata.txt.");

                return "";
            }

            std::string firstline;
            std::getline(file, firstline);

            file.close();

            return firstline;
        }
    };
}

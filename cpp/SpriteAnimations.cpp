#include "led-matrix.h"

#include <unistd.h>
#include <math.h>
#include <stdio.h>
#include <signal.h>

#include "constants.h"

#include "spritesheets.hpp"
#include "SpriteAnimations.hpp"

using rgb_matrix::Canvas;
using rgb_matrix::RGBMatrix;

namespace Sprites
{
    SpriteAnimations::SpriteAnimations(Canvas *canvas, volatile bool &interrupt_flag) : canvas_(canvas), interrupt_received(interrupt_flag) {}

    void SpriteAnimations::ScrollSpritesheet(Spritesheet spritesheet, int fps = 12)
    {
        int spriteIndex = 0;

        int maxSpriteWidth = spritesheet.width;
        int maxSpriteHeight = spritesheet.height;
        // int numSheets = spritesheet.numSheets;

        std::vector<int> colors = spritesheet.colors;

        int offsetX = 0 - maxSpriteWidth;
        int offsetY = (SCREEN_HEIGHT - maxSpriteHeight) / 2;

        while (!interrupt_received)
        {
            canvas()->Clear();

            std::vector<std::vector<int>> currentSprite = spritesheet.pixelData[spriteIndex];
            spriteIndex = (spriteIndex + 1) % spritesheet.pixelData.size();

            for (std::size_t y = 0; y < currentSprite.size(); ++y)
            {
                for (std::size_t x = 0; x < currentSprite[y].size(); ++x)
                {
                    int colorIndex = currentSprite[y][x];

                    if (colorIndex == -1)
                        continue;

                    int color = colors[colorIndex];

                    if (color == 0)
                        continue;

                    int red = (color >> 16) & 0xff;
                    int green = (color >> 8) & 0xff;
                    int blue = color & 0xff;

                    canvas()->SetPixel(x + offsetX, y + offsetY, red, green, blue);
                }
            }

            offsetX = offsetX + 3;
            if (offsetX > (SCREEN_WIDTH + maxSpriteWidth))
            {
                offsetX = 0 - maxSpriteWidth;
            }

            usleep((1000 / fps) * 1000);
        }
    };

    void SpriteAnimations::ShowAnimatedSpritesheet(Spritesheet spritesheet)
    {
    }
};

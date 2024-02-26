#pragma once

#include <vector>
#include <string>

struct Spritesheet
{
    int width;
    int height;
    int numSheets;
    std::vector<int> colors;
    std::vector<std::vector<std::vector<int>>> pixelData;
};

namespace Spritesheets
{
    extern Spritesheet KirbyWalking;
    extern Spritesheet KirbyTumbling;
    extern Spritesheet KirbyRunning;
}

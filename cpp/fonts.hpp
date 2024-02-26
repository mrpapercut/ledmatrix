#pragma once

#include <vector>
#include <string>
#include <unordered_map>

struct Font
{
    std::string name;
    std::unordered_map<char, std::vector<int>> characters;
};

namespace Fonts
{
    extern Font DefaultFont;
    extern Font SMWFont;
}

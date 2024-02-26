#pragma once

#include "led-matrix.h"
#include <vector>
#include "spritesheets.hpp"

using rgb_matrix::Canvas;

namespace Text
{
    class Text
    {
    protected:
        Text(Canvas *canvas, volatile bool &interrupt_flag);
        inline Canvas *canvas() { return canvas_; }

    public:
        virtual ~Text() {}
        virtual void Run() = 0;
        void ScrollText(Font font, const char *message);
        void ShowText(Font font, const char *message);

    private:
        Canvas *const canvas_;
        volatile bool &interrupt_received;
        Spritesheet ConvertFontToPixels(std::vector<std::vector<int>> fontData);
    };
};

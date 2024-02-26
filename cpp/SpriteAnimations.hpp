#pragma once

#include "led-matrix.h"
#include <vector>

using rgb_matrix::Canvas;

namespace Sprites
{

    class SpriteAnimations
    {
    protected:
        SpriteAnimations(Canvas *canvas, volatile bool &interrupt_flag);
        inline Canvas *canvas() { return canvas_; }

    public:
        virtual ~SpriteAnimations() {}
        virtual void Run() = 0;
        void ScrollSpritesheet(Spritesheet spritesheet, int fps);
        void ShowAnimatedSpritesheet(Spritesheet spritesheet);

    private:
        Canvas *const canvas_;
        volatile bool &interrupt_received;
    };
}

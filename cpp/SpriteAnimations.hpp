#pragma once

#include "led-matrix.h"
#include <vector>

#include "Spritesheet.hpp"

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

    class KirbyWalking : public SpriteAnimations
    {
    public:
        Spritesheet spriteData;

        KirbyWalking(Canvas *canvas, volatile bool &interrupt_flag)
            : SpriteAnimations(canvas, interrupt_flag), spriteData((Spritesheets::KirbyWalking)) {}

        void Run() override
        {
            int fps = 10;
            ScrollSpritesheet(spriteData, fps);
        }
    };

    class KirbyTumbling : public SpriteAnimations
    {
    public:
        Spritesheet spriteData;

        KirbyTumbling(Canvas *canvas, volatile bool &interrupt_flag)
            : SpriteAnimations(canvas, interrupt_flag), spriteData((Spritesheets::KirbyTumbling)) {}

        void Run() override
        {
            int fps = 15;
            ScrollSpritesheet(spriteData, fps);
        }
    };

    class KirbyRunning : public SpriteAnimations
    {
    public:
        Spritesheet spriteData;

        KirbyRunning(Canvas *canvas, volatile bool &interrupt_flag)
            : SpriteAnimations(canvas, interrupt_flag), spriteData((Spritesheets::KirbyRunning)) {}

        void Run() override
        {
            int fps = 15;
            ScrollSpritesheet(spriteData, fps);
        }
    };
}

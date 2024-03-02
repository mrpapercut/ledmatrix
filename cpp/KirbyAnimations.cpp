#include "Spritesheet.hpp"
#include "spritesheets.hpp"
#include "SpriteAnimations.hpp"

using rgb_matrix::Canvas;
using rgb_matrix::RGBMatrix;

namespace Sprites
{
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
            ShowAnimatedSpritesheet(spriteData);
        }
    };
}

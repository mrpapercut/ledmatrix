#include "led-matrix.h"

#include <unistd.h>
#include <math.h>
#include <memory>
#include <signal.h>
#include <stdio.h>

#include "constants.h"

#include "spritesheets.hpp"
#include "SpriteAnimations.hpp"

#include "Text.hpp"

using rgb_matrix::Canvas;
using rgb_matrix::RGBMatrix;

volatile bool interrupt_received = false;
static void InterruptHandler(int signo)
{
    interrupt_received = true;
};

static RGBMatrix::Options getDefaultOptions()
{
    RGBMatrix::Options defaults;

    defaults.hardware_mapping = "regular"; // or e.g. "adafruit-hat"
    defaults.cols = SCREEN_WIDTH;
    defaults.rows = SCREEN_HEIGHT;
    defaults.brightness = DEFAULT_BRIGHTNESS;
    defaults.disable_hardware_pulsing = true;

    return defaults;
};

static int usage(const char *progname)
{
    fprintf(stderr, "usage: %s <options> -M <mode> -T <type>\n",
            progname);
    fprintf(stderr, "Options:\n");
    fprintf(stderr,
            "\t-M <mode>: Required\n"
            "\t\tAvailable modes:\n"
            "\t\t0  - animate\n"
            "\n"
            "\t-T <type>: Required\n"
            "\t\tAvailable types:\n"
            "\t\t0  - Kirby walk\n"
            "\t\t1  - Kirby tumble\n"
            "\t\t2  - Kirby run\n"
            "\t\t3  - Show 'Hello, world!'\n"
            "\t\t4  - Show first line from inputdata.txt\n");
    return 1;
}

int main(int argc, char *argv[])
{
    RGBMatrix::Options defaults = getDefaultOptions();

    Canvas *canvas = RGBMatrix::CreateFromFlags(&argc, &argv, &defaults);
    if (canvas == NULL)
        return 1;

    int mode = -1;
    int type = -1;

    int opt;
    while ((opt = getopt(argc, argv, "m:M:t:T:")) != -1)
    {
        switch (opt)
        {
        case 'm':
        case 'M':
            mode = atoi(optarg);
            break;

        case 't':
        case 'T':
            type = atoi(optarg);
            break;

        default:
            return usage(argv[0]);
            break;
        }
    }

    if (mode != 0)
    {
        return usage(argv[0]);
    }

    signal(SIGTERM, InterruptHandler);
    signal(SIGINT, InterruptHandler);

    std::unique_ptr<Sprites::SpriteAnimations> animation = NULL;
    std::unique_ptr<Text::Text> textdisplay = NULL;

    switch (type)
    {
    case 0:
        animation = std::make_unique<Sprites::KirbyWalking>(canvas, interrupt_received);
        break;
    case 1:
        animation = std::make_unique<Sprites::KirbyTumbling>(canvas, interrupt_received);
        break;
    case 2:
        animation = std::make_unique<Sprites::KirbyRunning>(canvas, interrupt_received);
        break;
    case 3:
        textdisplay = std::make_unique<Text::ShowHelloWorld>(canvas, interrupt_received);
        break;
    case 4:
        textdisplay = std::make_unique<Text::ShowRandomLineFromFile>(canvas, interrupt_received);
        break;
    default:
        return usage(argv[0]);
    }

    if (animation != NULL)
    {
        animation->Run();
    }
    else if (textdisplay != NULL)
    {
        textdisplay->Run();
    }

    // Animation finished. Shut down the RGB matrix.
    canvas->Clear();
    delete canvas;

    return 0;
}

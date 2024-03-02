#pragma once

#include "led-matrix.h"
#include <vector>
#include <iostream>
#include <fstream>

#include "fonts.hpp"

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
    };

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
};

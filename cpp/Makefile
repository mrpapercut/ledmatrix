CFLAGS = -Wall -O3 -g -Wextra -Wno-unused-parameter
CXXFLAGS = $(CFLAGS)

BUILD_DIR := build
DIST_DIR := dist

OBJECTS := $(addprefix $(BUILD_DIR)/, spritesheets.o fonts.o Spritesheet.o Text.o SpriteAnimations.o LedMatrix.o)
BINARIES := $(addprefix $(DIST_DIR)/, LedMatrix)

RGB_LIB_DISTRIBUTION = ..
RGB_INCDIR = $(RGB_LIB_DISTRIBUTION)/include
RGB_LIBDIR = $(RGB_LIB_DISTRIBUTION)/lib
RGB_LIBRARY_NAME = rgbmatrix
RGB_LIBRARY = $(RGB_LIBDIR)/lib$(RGB_LIBRARY_NAME).a
LDFLAGS += -L$(RGB_LIBDIR) -l$(RGB_LIBRARY_NAME) -lrt -lm -lpthread

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

all: $(DIST_DIR) $(BINARIES)

$(RGB_LIBRARY):
	$(MAKE) -C $(RGB_LIBDIR)

$(DIST_DIR)/LedMatrix: $(OBJECTS) $(RGB_LIBRARY)
	$(CXX) $^ -o $@ $(LDFLAGS)

$(BUILD_DIR)/%.o: %.cpp | $(BUILD_DIR)
	$(CXX) -I$(RGB_INCDIR) $(CXXFLAGS) -c -o $@ $<

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)

.PHONY: all clean

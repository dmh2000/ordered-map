CXX = g++
CXXFLAGS = -std=c++11 -Wall -Wextra -pedantic

TARGET = orderedmap
SRCS = orderedmap.cpp

$(TARGET): $(SRCS)
	$(CXX) $(CXXFLAGS) -o $(TARGET) $(SRCS)

.PHONY: clean

clean:
	rm -f $(TARGET)

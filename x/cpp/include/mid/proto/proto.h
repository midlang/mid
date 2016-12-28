
#include <stdint.h>

namespace mid {
namespace proto {

	class Writer {
		virtual int write(const char* buf, size_t size) = 0;
		virtual int writeByte(char b) = 0;
	};

	class Reader {
		virtual int read(char* buf, size_t size) = 0;
		virtual int readByte(char* buf) = 0;
	};

} // namespace proto
} // namespace mid

#!/bin/bash

kinds='
default
beans
'

for kind in $kinds
do
midc \
	./demo.mid \
	-Ogo=generated/go_$kind \
	-Ocpp=generated/cpp_$kind \
	-Eautogen_decl="// NOTE: generated file, DON'T edit!!" \
	-Ecpp:unordered_map \
	-Xcodec \
	-K $kind \
	--log=info
done

#-Euse_fixed_encode \

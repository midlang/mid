#!/bin/bash

kinds='
default
beans
'

for kind in $kinds
do
midc \
	-Ogo=generated/go_$kind \
	-Ocpp=generated/cpp_$kind \
	-Eautogen_decl="// NOTE: generated file, DON'T edit!!" \
	-Ecpp:unordered_map \
	-Xcodec \
	-K $kind \
	--log=debug \
	./demo.mid
done

#-Euse_fixed_encode \

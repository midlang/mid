#!/bin/bash

midc \
	-I ./testdata/ \
	-Ogo=./testdata/generated/go \
	-Tgo=./testdata/templates/go/ \
	-vi

const ecojiPromise = WebAssembly.instantiateStreaming(
    fetch("ecoji.wasm")
)
    .then((result) => {
        const {
            ecoji_encode_bytes,
            ecoji_decode_bytes,
            ecoji_allocate,
            ecoji_deallocate,
            ecoji_encode_len,
            ecoji_decode_len,
            memory
        } = result.instance.exports;

        function allocateAndCopy(input, outputLen) {

            // allocate memory within wasm memory space for the input data
            const inputStart = ecoji_allocate(input.length);
            console.log(`input start : ${inputStart} input.length:${input.length}`);

            if (inputStart < 0) {
                throw "unable to allocate memory";
            }

            // allocate memory within wasm memory space for the output data
            const outputStart = ecoji_allocate(outputLen);
            if (outputStart < 0) {
                ecoji_deallocate(inputStart, input.length);
                throw "unable to allocate memory";
            }
            console.log(`outputStart : ${outputStart}`);

            // copy the input data from javascript memory to wasm memory
            const memView = new Uint8Array(memory.buffer);
            for (let i = 0; i < input.length; i++) {
                memView[inputStart + i] = input[i];
            }

            return {inputAddr: inputStart, inputLen: input.length, outputAddr: outputStart, outputLen: outputLen};
        }

        function deallocateAndCopy(allocations, actualOutputLen) {
            // copy the data from wasm memory to javascript memory
            const outputCopy = actualOutputLen === 0 ? "" :
                new TextDecoder().decode(new Uint8Array(memory.buffer, allocations.outputAddr, actualOutputLen));

            //dealloc after copying as it may zero wasm memory
            ecoji_deallocate(allocations.inputAddr, allocations.inputLen);
            ecoji_deallocate(allocations.outputAddr, allocations.outputLen);

            return outputCopy;
        }

        function ecojiEncode(inputStr) {
            console.log("ecojiEncode:", inputStr);

            // convert input string to a utf8 byte array
            const input = new TextEncoder().encode(inputStr);

            // compute the maximum size that will be needed for output data
            const maxOutputLen = ecoji_encode_len(input.length);
            //console.log(`maxOutputLen : ${maxOutputLen}`);

            // allocate memory in wasm and copy input data into wasm memory
            const allocations = allocateAndCopy(input, maxOutputLen);

            // call the ecoji encoding function that is wasm code
            const outputLen = ecoji_encode_bytes(allocations.inputAddr, allocations.inputLen, allocations.outputAddr, allocations.outputAddr);

            console.log("outputLen: " + outputLen);

            // check for errors encoding
            if (outputLen < 0) {
                return "ERROR:Failed to encode input!";
            }

            return deallocateAndCopy(allocations, outputLen);
        }

        function ecojiDecode(inputStr) {
            console.log("ecojiDecode: " + inputStr);

            // encode input string as utf8 array
            const input = new TextEncoder().encode(inputStr);

            // compute the maximum size that will be needed for output data
            const maxOutputLen = ecoji_decode_len(input.length);

            // allocate memory in wasm and copy input data into wasm memory
            const allocations = allocateAndCopy(input, maxOutputLen);

            const outputLen = ecoji_decode_bytes(allocations.inputAddr, allocations.inputLen, allocations.outputAddr, allocations.outputLen);

            console.log("outputLen: " + outputLen);

            if (outputLen < 0) {
                return "🤨 It seems the input was not Ecoji™ encoded";
            }

            return deallocateAndCopy(allocations, outputLen);
        }

        return {encode: ecojiEncode, decode: ecojiDecode};
    });


async function ecojiDecode(inputStr) {
    const ecoji = await ecojiPromise;
    return ecoji.decode(inputStr);
}

async function ecojiEncode(inputStr) {
    const ecoji = await ecojiPromise;
    return ecoji.encode(inputStr);
}

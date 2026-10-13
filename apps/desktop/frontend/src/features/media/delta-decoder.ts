const MAGIC = "NDF1";
const HEADER_SIZE = 20;
const MODE_SOLID = 1;
const MODE_XOR = 2;

export type DeltaFrame = {
  width: number;
  height: number;
  tileSize: number;
  keyframe: boolean;
  tiles: number;
};

/** Applies one Nodex delta packet to an existing canvas. */
export async function applyDeltaPacket(
  context: CanvasRenderingContext2D,
  packet: ArrayBuffer,
  previousWidth: number,
  previousHeight: number,
): Promise<DeltaFrame> {
  const bytes = new Uint8Array(packet);
  const view = new DataView(packet);
  const magic = new TextDecoder().decode(bytes.subarray(0, 4));
  if (magic !== MAGIC || bytes[4] !== 1)
    throw new Error("Invalid Nodex delta packet");

  const keyframe = bytes[5] === 1;
  const width = view.getUint32(6);
  const height = view.getUint32(10);
  const tileSize = view.getUint16(14);
  const tileCount = view.getUint32(17);
  if (!keyframe && (width !== previousWidth || height !== previousHeight)) {
    throw new Error("Delta frame dimensions do not match the canvas");
  }
  if (
    keyframe &&
    (context.canvas.width !== width || context.canvas.height !== height)
  ) {
    context.canvas.width = width;
    context.canvas.height = height;
  }

  let offset = HEADER_SIZE;
  for (let index = 0; index < tileCount; index += 1) {
    const x = view.getUint16(offset);
    const y = view.getUint16(offset + 2);
    const tileWidth = view.getUint16(offset + 4);
    const tileHeight = view.getUint16(offset + 6);
    const mode = bytes[offset + 8];
    const payloadLength = view.getUint32(offset + 9);
    offset += 13;
    const compressed = bytes.subarray(offset, offset + payloadLength);
    offset += payloadLength;
    const data = new Uint8Array(await decompress(compressed));
    const pixels =
      mode === MODE_SOLID
        ? expandSolid(data, tileWidth * tileHeight)
        : mode === MODE_XOR
          ? xorWithCanvas(context, data, x, y, tileWidth, tileHeight)
          : data;
    context.putImageData(
      new ImageData(new Uint8ClampedArray(pixels), tileWidth, tileHeight),
      x,
      y,
    );
  }
  return { width, height, tileSize, keyframe, tiles: tileCount };
}

async function decompress(data: Uint8Array): Promise<ArrayBuffer> {
  if (typeof DecompressionStream === "undefined")
    throw new Error("DecompressionStream is unavailable");
  const owned = new Uint8Array(data);
  const stream = new Blob([owned.buffer])
    .stream()
    .pipeThrough(new DecompressionStream("deflate"));
  return streamToArrayBuffer(stream);
}

async function streamToArrayBuffer(
  stream: ReadableStream<Uint8Array>,
): Promise<ArrayBuffer> {
  const reader = stream.getReader();
  const chunks: Uint8Array[] = [];
  let length = 0;
  for (;;) {
    const part = await reader.read();
    if (part.done) break;
    chunks.push(part.value);
    length += part.value.length;
  }
  const result = new Uint8Array(length);
  let offset = 0;
  for (const chunk of chunks) {
    result.set(chunk, offset);
    offset += chunk.length;
  }
  return result.buffer;
}

function expandSolid(color: Uint8Array, pixels: number): Uint8Array {
  const result = new Uint8Array(pixels * 4);
  for (let i = 0; i < pixels; i += 1) result.set(color, i * 4);
  return result;
}

function xorWithCanvas(
  context: CanvasRenderingContext2D,
  data: Uint8Array,
  x: number,
  y: number,
  width: number,
  height: number,
): Uint8Array {
  const previous = context.getImageData(x, y, width, height).data;
  for (let i = 0; i < data.length; i += 1) data[i] ^= previous[i];
  return data;
}

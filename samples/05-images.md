# Images

Local images are automatically embedded as base64 data URIs in the PDF.

## Local Image

> To test this sample, place a `sample.png` file in this directory.

![Sample Image](sample.png)

## Remote Images

Remote images (http/https) are kept as-is and require network access during rendering.

## How It Works

1. goldmark converts `![alt](path)` to `<img src="path">`
2. unimd2pdf detects local `src` paths
3. Reads the file and converts to `data:image/png;base64,...`
4. chromedp renders the embedded image into the PDF

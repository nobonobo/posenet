build:
	esbuild index.js --format=esm --global-name=posenet --bundle --platform=node --outfile=../posenet.js

const express = require("express");
//const next = require("next");
//const compression = require("compression");

const port = parseInt(process.env.PORT, 10) || 4440;
//const prod = process.env.NODE_ENV === "production";
//const app = next({ dir: "./src", dev: !prod });
//const handle = app.getRequestHandler();
const FS = require("fs");

var path = require('path');

function rewrite(apiPath, extensions, req, res, next) {
	const url = req.url.split("?");
	const b = `${url[0]}/${req.method}`;
	let js = false;
	extensions.some(ext => {
		const resourcePath = `${b}.${ext}`;
		const filePath = path.posix.join(apiPath, resourcePath);
		if (FS.existsSync(filePath)) {
			if (ext === "script") {
				console.log(filePath);
				const rew = require(filePath);
				js = true;
				return rew.rewrite(req, res);
			}
			req.method = "GET";
			req.url = resourcePath;
			return true;
		}
		return false;
	});
	if (!js) next();
	console.log(req.method, req.originalUrl, "→", req.url);
}

//app.prepare().then(() => {
	const server = express();
	/*if (prod) {
		server.use(compression());
	} else {*/
		var appDir = path.dirname(require.main.filename);
		const apiPath = appDir.concat(`/api`);
		const extensions = ["js", "json", "jpg", "html", "pdf", "png", "txt", "script"];
		server.use("/api", (req, res, next) => rewrite(apiPath, extensions, req, res, next));
		server.use("/api", express.static(apiPath));
	//}

	//server.get("*", (req, res) => handle(req, res));

	server.listen(port, err => {
		if (err) throw err;
		console.log(`> Ready on http://localhost:${port}`);
	});
//});

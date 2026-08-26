async function api(path, body) {
	const options = {method: "GET", headers: {"Content-Type": "application/json"}};
	if (body !== undefined) {
		options.method = "POST";
		options.body = JSON.stringify(body);
	}
	const res = await fetch(path, options);
	const data = await res.json();
	if (!data.ok) {
		throw new Error(data.error || "request failed");
	}
	return data.data;
}

function el(selector) {
	return document.querySelector(selector);
}

function renderStatus(data) {
	const box = el("#status");
	if (!box) {
		return;
	}
	const water = data.water || {};
	const brush = data.brush || {};
	const roof = data.roof || {};
	const dry = data.dry || {};
	const chem = data.chem || {};
	const conv = data.conv || {};
	const cycle = data.cycle || {};
	box.innerHTML =
		"<ul>" +
		"<li>工艺阶段: " + cycle.stage + "</li>" +
		"<li>列车: " + (data.position && data.position.train_id || "-") + "</li>" +
		"<li>水系统: " + water.state + "</li>" +
		"<li>侧刷: " + (brush.lowered ? "下放" : "抬起") + "</li>" +
		"<li>顶刷: " + (roof.lowered ? "下放" : "抬起") + "</li>" +
		"<li>风干机: " + (dry.running ? "运行" : "停止") + "</li>" +
		"<li>药剂浓度: " + (chem.dose_ok ? "正常" : "偏高") + "</li>" +
		"<li>输送带: " + (conv.running ? "运行" : "停止") + "</li>" +
		"</ul>";
}

function renderSegments(data) {
	const box = el("#segments");
	if (!box) {
		return;
	}
	const rows = (data.segments || []).map(function (seg) {
		return "<tr><td>" + seg.kind + "</td><td>" + seg.start_mm + "</td><td>" + seg.end_mm + "</td><td>" + (seg.contains_front ? "是" : "否") + "</td></tr>";
	}).join("");
	box.innerHTML = "<table><thead><tr><th>段位</th><th>起</th><th>止</th><th>含车头</th></tr></thead><tbody>" + rows + "</tbody></table>";
}

async function loadAudit() {
	const table = el("#audit");
	if (!table) {
		return;
	}
	const data = await api("/api/audit");
	const rows = (data.records || []).map(function (rec) {
		return "<tr><td>" + rec.created_at + "</td><td>" + rec.action + "</td><td>" + rec.train_id + "</td><td>" + rec.detail + "</td></tr>";
	}).join("");
	table.querySelector("tbody").innerHTML = rows;
}

async function refresh() {
	try {
		const data = await api("/api/state");
		renderStatus(data);
		renderSegments(data);
		const brushState = el("#brush-state");
		if (brushState) {
			brushState.textContent = "刷组: " + (data.brush.active_group.name || "-") + " 下放=" + data.brush.lowered;
		}
		const waterState = el("#water-state");
		if (waterState) {
			waterState.textContent = "状态=" + data.water.state + " 增益=" + data.water.gain_mpa;
		}
		const chemState = el("#chem-state");
		if (chemState) {
			chemState.textContent = "告警=" + data.chem.alarm + " 阀锁=" + data.chem.valve_latched + " 剂量=" + data.chem.dose_ml;
		}
		const cycleState = el("#cycle-state");
		if (cycleState) {
			cycleState.textContent = "阶段=" + data.cycle.stage + " 列车=" + data.cycle.train_id;
		}
		await loadAudit();
	} catch (err) {
		console.error(err);
	}
}

function bind(selector, path, bodyFactory, resultSelector) {
	const button = el(selector);
	if (!button) {
		return;
	}
	button.addEventListener("click", async function () {
		try {
			const body = bodyFactory ? bodyFactory() : {};
			const data = await api(path, body);
			if (resultSelector) {
				el(resultSelector).textContent = JSON.stringify(data);
			}
			await refresh();
		} catch (err) {
			if (resultSelector) {
				el(resultSelector).textContent = String(err.message || err);
			}
		}
	});
}

bind("#lower-btn", "/api/brush/lower", null, "#brush-result");
bind("#raise-btn", "/api/brush/raise", null, "#brush-result");
bind("#retract-btn", "/api/brush/retract", null, "#brush-result");
bind("#roof-lower-btn", "/api/roof/lower", null, "#roof-result");
bind("#roof-raise-btn", "/api/roof/raise", null, "#roof-result");
bind("#roof-brush-btn", "/api/roof/brush", null, "#roof-result");
bind("#water-start-btn", "/api/water/start", null, "#water-result");
bind("#water-rinse-btn", "/api/water/rinse", null, "#water-result");
bind("#water-stop-btn", "/api/water/stop", null, "#water-result");
bind("#water-drain-btn", "/api/water/drain", null, "#water-result");
bind("#spray-btn", "/api/chem/spray", function () { return {duration_ms: 15000}; }, "#chem-result");
bind("#alarm-btn", "/api/chem/alarm", null, "#chem-result");
bind("#alarm-clear-btn", "/api/chem/alarm/clear", null, "#chem-result");
bind("#dry-start-btn", "/api/dry/start", null, "#dry-result");
bind("#dry-stop-btn", "/api/dry/stop", null, "#dry-result");
bind("#complete-btn", "/api/plan/complete", null, "#cycle-result");
bind("#stop-btn", "/api/plan/stop", null, "#cycle-result");

const washForm = el("#wash-form");
if (washForm) {
	washForm.addEventListener("submit", async function (event) {
		event.preventDefault();
		const form = new FormData(washForm);
		try {
			const data = await api("/api/plan/wash", {
				train_id: form.get("train_id"),
				type: form.get("type"),
				length_mm: Number(form.get("length_mm")),
				front_mm: Number(form.get("front_mm")),
				zero_mm: Number(form.get("zero_mm"))
			});
			el("#wash-result").textContent = JSON.stringify(data);
			await refresh();
		} catch (err) {
			el("#wash-result").textContent = String(err.message || err);
		}
	});
}

const typeForm = el("#type-form");
if (typeForm) {
	typeForm.addEventListener("submit", async function (event) {
		event.preventDefault();
		const form = new FormData(typeForm);
		try {
			const data = await api("/api/entry/type", {type: form.get("type")});
			el("#type-result").textContent = JSON.stringify(data);
			await refresh();
		} catch (err) {
			el("#type-result").textContent = String(err.message || err);
		}
	});
}

const gainForm = el("#gain-form");
if (gainForm) {
	gainForm.addEventListener("submit", async function (event) {
		event.preventDefault();
		const form = new FormData(gainForm);
		try {
			const data = await api("/api/water/recalibrate", {gain_mpa: Number(form.get("gain_mpa"))});
			el("#gain-result").textContent = JSON.stringify(data);
			await refresh();
		} catch (err) {
			el("#gain-result").textContent = String(err.message || err);
		}
	});
}

refresh();
setInterval(refresh, 3000);

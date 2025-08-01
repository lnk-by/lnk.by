const params = new URLSearchParams(window.location.search);
const styleHref = params.get("style");
const confUrl = params.get("conf") || window.location.pathname.split("/").pop().replace(/\.\w+$/, ".json");
if (styleHref) {
    document.getElementById("page-style").href = styleHref;
}

function applyProperties(target, properties) {
    for (const [key, value] of Object.entries(properties)) {
        if (typeof value === "object" && value !== null && key in target) {
            applyProperties(target[key], value);  // recurse into nested object
        } else {
            target[key] = value;
        }
    }
}

// Load and apply config
fetch(confUrl)
    .then(res => res.json())
    .then(config => {
        for (const [elementId, properties] of Object.entries(config)) {
            const el = document.getElementById(elementId);
            if (!el) continue;
            applyProperties(el, properties);
        }

        
        const auto = config.auto
        if (auto && auto.href) {
            const delay = auto.delay || 0
            if (delay) {
                setTimeout(() => {
                    window.location.href = auto.href;
                }, auto.delay);
            } else {
                window.location.href = auto.href;
            }
        }
    })
    .catch(err => console.error("Failed to load config:", err));

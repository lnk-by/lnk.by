async function initConfiguration() {
  const res = await fetch('env-config.json');
  const config = await res.json();

  const selector = document.getElementById("apiSelector");

  selector.appendChild(new Option("localhost", "http://localhost:8080"));

  for (const [envName, envData] of Object.entries(config)) {
    const url = `https://${envData.apiId}.execute-api.${envData.region}.amazonaws.com`;
    selector.appendChild(new Option(envName, url));
  }
  return config;
}

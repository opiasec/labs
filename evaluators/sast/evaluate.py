import subprocess
import json
import os
import sys

RESULT_JSON = ""
SCAN_PATH = ""

SEVERITY_SCORES = {
    "ERROR": -10,
    "WARNING": -5,
    "INFO": -2,
}

def run_semgrep():
    os.makedirs("/opt/lab/results/sast", exist_ok=True)
    try:
        subprocess.run([
            "semgrep", "scan",
            "--quiet",
            "--config=auto",
            SCAN_PATH,
            "--json",
            "--output", RESULT_JSON
        ], check=True)
        return 0
    except subprocess.CalledProcessError as e:
        return e.returncode

def calculate_score(results):
    score = 100
    for finding in results:
        severity = finding.get("extra", {}).get("severity", "").upper()
        score += SEVERITY_SCORES.get(severity, 0)
    return max(score, 0)

def format_findings(findings):
    formatted = []
    for f in findings:
        formatted.append({
            "file": f.get("path"),
            "line": f.get("start", {}).get("line"),
            "severity": f.get("extra", {}).get("severity"),
            "id": f.get("check_id"),
            "message": f.get("extra", {}).get("message")
        })
    return formatted

def main():
    exit_code = run_semgrep()

    if not os.path.exists(RESULT_JSON):
        print(json.dumps({
            "score": 0,
            "message": "Semgrep falhou e não gerou resultado."
        }))
        return

    with open(RESULT_JSON, "r") as f:
        data = json.load(f)

    if exit_code != 0:
        message = data.get("errors", [{}])[0].get("message", "Erro desconhecido executando o Semgrep.")
        print(json.dumps({
            "score": 0,
            "message": message
        }))
        return

    results = data.get("results", [])
    if not results:
        print(json.dumps({
            "score": 100,
            "message": "Nenhuma vulnerabilidade encontrada."
        }))
        return

    score = calculate_score(results)
    findings = format_findings(results)

    findings_json_string = json.dumps(findings)

    print(json.dumps({
        "score": score,
        "message": findings_json_string
    }))

if __name__ == "__main__":
    SCAN_PATH = sys.argv[1]
    RESULT_JSON = sys.argv[2]
    main()

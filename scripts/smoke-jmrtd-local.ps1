# Local JMRTD smoke — exit-code only (do not treat stderr as failure)
#
# Root cause of 20 Jul false-fail: PowerShell `$ErrorActionPreference = 'Stop'`
# turns any native-command stderr write into terminating NativeCommandError,
# even when java exits 0. Diagnostic `PACE_FAIL_EXCEPTION_CLASS=…` is on stderr
# by design. Suite cells must assert `$LASTEXITCODE` / process returncode only.
#
# Usage (from harness repo root):
#   powershell -File scripts/smoke-jmrtd-local.ps1

$ErrorActionPreference = "Continue"  # never 'Stop' around java
$root = Split-Path -Parent $PSScriptRoot
$jar = Join-Path $root "drivers\jmrtd\target\jmrtd-tc-ac-01-0.1.0.jar"
if (-not (Test-Path $jar)) {
    Write-Error "Missing $jar — build fat jar first"
    exit 2
}
$log = Join-Path $root ("logs\smoke-jmrtd-0.8.6-" + (Get-Date -Format "yyyyMMddTHHmmss") + "Z")
New-Item -ItemType Directory -Force -Path $log | Out-Null
Set-Location $root

function Invoke-JavaCell([string]$main, [string]$label) {
    Write-Host "=== $label ==="
    & java -cp $jar $main -profile profiles/pace-then-bac-downgrade.json -log-dir $log
    $rc = $LASTEXITCODE
    Write-Host "$label LASTEXITCODE=$rc"
    return $rc
}

$rc1 = Invoke-JavaCell "org.emrtd.harness.jmrtd.TcAc01Runner" "baseline"
$rc2 = Invoke-JavaCell "org.emrtd.harness.jmrtd.TcAc01MitigatedRunner" "mitigated"
# Assert *before* any Format-Table / listing — those can clobber $LASTEXITCODE on Windows PS.
$failed = ($rc1 -ne 0) -or ($rc2 -ne 0)
Write-Host "artifacts under $log (rc1=$rc1 rc2=$rc2)"
Get-ChildItem $log | Format-Table Name, Length | Out-Host
if ($failed) { exit 1 }
exit 0

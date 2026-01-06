param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("start","stop","status")]
    [string]$Action
)

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot '..') | Select-Object -ExpandProperty Path
$pidFile = Join-Path $repoRoot '.run_loop.pid'

function Get-ProcessAlive($id) {
    try { Get-Process -Id $id -ErrorAction Stop | Out-Null; return $true } catch { return $false }
}

switch ($Action) {
    "start" {
        if (Test-Path $pidFile) {
            $existing = Get-Content $pidFile -ErrorAction SilentlyContinue
            if ($existing -and (Get-ProcessAlive $existing)) {
                Write-Output "Loop already running with PID $existing"
                exit 0
            } else {
                Remove-Item $pidFile -ErrorAction SilentlyContinue
            }
        }
        $worker = Join-Path $PSScriptRoot 'loop_worker.ps1'
        $proc = Start-Process -FilePath pwsh -ArgumentList "-NoProfile","-ExecutionPolicy","Bypass","-File","`"$worker`"" -PassThru
        $proc.Id | Out-File -FilePath $pidFile -Encoding ascii
        Write-Output "Started loop worker with PID $($proc.Id)"
    }
    "stop" {
        if (Test-Path $pidFile) {
            $workerPid = (Get-Content $pidFile) -as [int]
            if ($workerPid -and (Get-ProcessAlive $workerPid)) {
                Stop-Process -Id $workerPid -Force
                Write-Output "Stopped process $workerPid"
            } else {
                Write-Output "No running process found for PID $workerPid"
            }
            Remove-Item $pidFile -ErrorAction SilentlyContinue
        } else {
            Write-Output "No PID file found; nothing to stop."
        }
    }
    "status" {
        if (Test-Path $pidFile) {
            $workerPid = (Get-Content $pidFile) -as [int]
            if ($workerPid -and (Get-ProcessAlive $workerPid)) {
                Write-Output "Loop running with PID $workerPid"
            } else {
                Write-Output "PID file found but process not running."
            }
        } else {
            Write-Output "Loop is not running."
        }
    }
}

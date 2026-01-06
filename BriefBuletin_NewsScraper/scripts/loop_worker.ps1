# Worker loop: runs main.py, commits/pushes if changes, sleeps 30 minutes
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot '..') | Select-Object -ExpandProperty Path
Set-Location $repoRoot

function Get-PythonCommand {
    # Try common python executables first
    $candidates = @('python', 'python3', 'py')
    foreach ($name in $candidates) {
        $cmd = Get-Command $name -ErrorAction SilentlyContinue
        if ($cmd) { return $cmd.Source }
    }

    # Fallback to where.exe which lists installed locations
    try {
        $where = (& where.exe python) -join "`n" 2>$null
        if ($where) {
            $first = $where -split "`n" | Select-Object -First 1
            if ($first) { return $first.Trim() }
        }
    } catch {}

    # Check common environment variables
    if ($env:PYTHON) { return $env:PYTHON }
    if ($env:PYTHON_HOME) { return (Join-Path $env:PYTHON_HOME 'python.exe') }

    return $null
}

$pythonCmd = Get-PythonCommand
if (-not $pythonCmd) {
    Write-Output "WARNING: No Python executable found in PATH. Set 'python' on PATH or set the PYTHON env variable."
}

while ($true) {
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-Output "[$timestamp] Running main.py..."

    if ($pythonCmd) {
        try {
            Write-Output "Using Python: $pythonCmd"
            & $pythonCmd .\main.py
            $exit = $LASTEXITCODE
            if ($exit -ne 0) { Write-Output "main.py exited with code $exit" }
        } catch {
            Write-Output "Error running main.py: $_"
        }
    } else {
        Write-Output "Error: No Python executable available. Will retry after sleep."
    }

    # check for git changes
    $changes = (& git status --porcelain) -join "`n"
    if ($changes) {
        Write-Output "Changes detected, committing..."
        & git add .
        $msg = "last update $((Get-Date).ToString('yyyy-MM-dd HH:mm:ss'))"
        & git commit -m $msg
        $commitExit = $LASTEXITCODE
        if ($commitExit -ne 0) { Write-Output "git commit returned $commitExit" }
        Write-Output "Pushing..."
        & git push
        $pushExit = $LASTEXITCODE
        if ($pushExit -ne 0) { Write-Output "git push returned $pushExit" }
    } else {
        Write-Output "No changes to commit."
    }

    Write-Output "Sleeping 30 minutes..."
    Start-Sleep -Seconds 1800 #600  # 1800 for few days
    # Re-evaluate python command each loop in case PATH changes
    if (-not $pythonCmd) { $pythonCmd = Get-PythonCommand }
}

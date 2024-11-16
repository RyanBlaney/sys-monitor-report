## Categories of System Data

#### CPU Usage

**Metrics**:
- [x] Overall CPU usage (%).
- [x] Per-core CPU usage (%).
- [x] Top applications by CPU usage (%).
- [ ] Largest spikes in CPU usage (e.g., over a threshold).

#### Memory (RAM) Usage

**Metrics**:
- [x] Total memory available and used (MB/GB).
- [x] Percentage of memory used.
- [x] Top applications by memory usage.
- [ ] Largest spikes in memory usage.

#### Swap Memory

**Metrics**:
- [x] Swap memory used and free (MB/GB).
- [ ] Top applications using swap memory (if possible).

#### Disk Usage

**Metrics**:
- [x] Total disk space available and used (MB/GB).
- [ ] Disk I/O read/write speed (MB/s).
- [x] Top applications with highest disk I/O activity.
- [ ] Spikes in disk read/write activity.

#### Network Usage

**Metrics**:
- [ ] Total data sent and received (MB/GB).
- [ ] Current upload and download speed (MB/s).
- [ ] Top applications by network usage.
- [ ] Applications causing network spikes (e.g., exceeding a threshold).

#### Systemd Processes

**Metrics**:
- [ ] Process ID, name, CPU, and memory usage.
- [ ] Status (running, sleeping, etc.).
- [ ] Priority and threads.

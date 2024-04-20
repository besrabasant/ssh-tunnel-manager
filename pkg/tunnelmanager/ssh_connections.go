package tunnelmanager

// SSH connection
type SSHConnections map[int]*ConnectionInfo

func (c SSHConnections) Filter(predicate func(*ConnectionInfo) bool) SSHConnections {
	filteredConns := make(SSHConnections, 0)

	for port, entry := range c {
		entry := entry
		if predicate(entry) {
			filteredConns[port] = entry
		}
	}

	return filteredConns
}

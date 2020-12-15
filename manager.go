// manager.go
package workflow

type Manager struct {
	wfs map[string]*WorkFlow
}

func (m Manager) len() int {
	return len(m.wfs)
}

func (m Manager) workFlow(id string) *WorkFlow {
	return m.wfs[id]
}
func (m Manager) workFlows() map[string]*WorkFlow {
	return m.wfs
}

func (m Manager) addWorkFlow(wf *WorkFlow) {
	m.wfs[wf.fId] = wf
}

func (m Manager) runWorkFlow(wf *WorkFlow) {

}

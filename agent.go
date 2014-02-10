package gosnmp

import (
	"code.google.com/p/biogo.llrb"
)

type TransactionProvider interface {
	// StartTxn should begin a transaction and return an opaque reference to the new transaction. If a transaction can't
	// be started, a nil value should be returned.
	StartTxn() interface{}
	// CommitTxn should commit the given transaction. If the transaction can't be committed for any reason, it should be
	// aborted and CommitTxn should return false. Otherwise, CommitTxn should return true.
	CommitTxn(interface{}) bool
	// AbortTxn should abort the given transaction, performing any rollback required.
	AbortTxn(interface{})
}

type Agent struct {
	snmpContext
	oidTree     llrb.Tree
	txnProvider TransactionProvider
}

func NewAgent(name string, maxTargets int, logger Logger, txnProvider TransactionProvider) *Agent {
	return NewAgentWithPort(name, maxTargets, 161, logger, txnProvider)
}

func NewAgentWithPort(name string, maxTargets int, port int, logger Logger, txnProvider TransactionProvider) *Agent {
	agent := new(Agent)
	agent.snmpContext.initContext(name, maxTargets, false, port, logger)
	agent.incomingRequestProcessor = agent
	agent.oidTree = llrb.Tree{}
	agent.txnProvider = txnProvider
	return agent
}

func (agent *Agent) processcommunityRequest(req *communityRequest) {
	resp := req.createResponse()
	txn := agent.txnProvider.StartTxn()
	if txn == nil {
		resp.errorVal = SnmpRequestErrorType_RESOURCE_UNAVAILABLE
		resp.errorIdx = 1
	}
	for _, requestVb := range req.varbinds {
		node := agent.lookupHandler(requestVb.GetOid())
		if node == nil {
			resp.AddVarbind(NewNoSuchObjectVarbind(requestVb.GetOid()))
			continue
		}
		switch req.pduType {
		case pduType_GET_REQUEST:
			responseVb, err := node.handler.Get(requestVb.GetOid(), txn)
			if err != nil {
				continue
			}
			resp.AddVarbind(responseVb)

		case pduType_SET_REQUEST:
			responseVb, err := node.handler.Set(requestVb, txn)
			if err != nil {
				continue
			}
			resp.AddVarbind(responseVb)

		}
	}
	agent.sendResponse(resp)
}

func (agent *Agent) lookupHandler(oid ObjectIdentifier) *oidTreeNode {
	tnode := agent.oidTree.Ceil(oidTreeLookup(oid))
	if tnode == nil {
		// This should only ever hit if no handlers have been added to this agent... Very much a corner case.
		agent.Errorf("------ Ctxt %s, YOU APPEAR TO HAVE NO HANDLERS BOUND", agent.name)
		return nil
	}
	node := tnode.(*oidTreeNode)
	if node.oid.MatchLength(oid) != len(node.oid) {
		// The node we looked up doesn't match the request OID. Note that it's ok for the request OID to be more
		// specific than the OID specified by the handler... in fact for all but the simplest requests, it's pretty much
		// guaranteed, where the request OID will specify a row in a table, or some extra information on top of the base
		// handler OID.
		return nil
	}
	return node
}

type oidHandler interface {
	Get(oid ObjectIdentifier, txn interface{}) (Varbind, error)
	Set(vb Varbind, txn interface{}) (Varbind, error)
}

type SingleVarOidHandler interface {
	oidHandler
}

func (agent *Agent) RegisterSingleVarOidHandler(oid ObjectIdentifier, handler SingleVarOidHandler) error {
	agent.oidTree.Insert(&oidTreeNode{oid, false, handler})
	return nil
}

type oidTreeNode struct {
	oid     ObjectIdentifier
	isMulti bool
	handler oidHandler
}

func (a *oidTreeNode) Compare(b llrb.Comparable) int {
	return a.oid.Compare(b.(*oidTreeNode).oid)
}

type oidTreeLookup ObjectIdentifier

func (a oidTreeLookup) Compare(b llrb.Comparable) int {
	return ObjectIdentifier(a).Compare(b.(*oidTreeNode).oid)
}

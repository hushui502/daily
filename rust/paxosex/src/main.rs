
#![allow(unused)]
#![allow(non_camel_case_types)]

fn main() {
    println!("Hello, world!");
}

mod processor {
    use crate::roles::*;

    pub struct Processor {
        roles: Roles,
        others: Vec<RemoteProcessor>,
    }

    pub struct Roles {
        client: Option<Client>,
        acceptor: Option<Acceptor>,
        proposor: Option<Proposer>,
        learner: Option<Learner>,
        leader: Option<Leader>,
    }

    struct Id(u128);

    pub struct RemoteProcessor {
        id: Id,
        client: bool,
        acceptor: bool,
        proposor: bool,
        learner: bool,
        leader: bool,
    }
}

mod roles {
    pub struct Value(u32);
    pub struct ProposalNumber(u128);

    pub struct Client;

    pub struct Acceptor {
        pub (crate) highest_proposal_number: ProposalNumber,
        pub (crate) last_accepted: Option<(ProposalNumber, Value)>,
    }

    pub struct Proposer {
        pub (crate) next_proposal_number: ProposalNumber,
    }

    pub struct Learner;

    pub struct Leader;
}

mod basic_paxos {
    use crate::roles::*;

    enum Message {
        Prepare(PrepareMessage)
    }

    struct PrepareMessage(ProposalNumber);

    enum BasicPaxosRound {
        Phase1(BasicPaxosRoundPhase1),
        Phase2(BasicPaxosRoundPhase2),
    }

    enum BasicPaxosRoundPhase1 {
        Phase1A_Prepare,
        Phase1B_Promise,
    }

    enum BasicPaxosRoundPhase2 {
        Phase2A_Accept,
        Phase2B_Accepted,
    }

    trait BasicProposer {
        fn create_prepare_msg(&mut self) -> Message;
    }

    impl BasicProposer for Proposer {
        fn create_prepare_msg(&mut self) -> Message {
            let number = self.next_prepare_number;
            self.next_prepare_number += 1;
            PrepareMsg(PrepareMessage(number))
        }
    }

    trait BasicAcceptor {
        fn handle_prepare_msg(&mut self, msg: PrepareMsg) -> HandlePrepareResult {
            panic!()
        }
    }


    enum HandlePrepareResult {
        Promise(Promise),
    }

    struct Promise {

    }
}

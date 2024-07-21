# Terminology
- `State Chart` \
A state chart, also known as a state diagram or state machine diagram,
 is a graphical representation of a finite state machine (FSM). 
 It depicts the various states an object can be in and the transitions 
 between these states based on events and conditions. 
 State charts are used in computer science and systems engineering 
 to model and understand the behavior of systems 
 that can be in one of a finite number of states at any given time.

- `Controller` \
In a state chart, the controller plays a critical role in 
guiding the system's progression through various states. 
It actively manages the execution of actions, 
assesses the validity of guards, 
initiates transitions to new states, 
and responds to any encountered errors.

- `Reconciler` \
In the context of state charts or state machines, 
a "reconciler" is a mechanism or component that is responsible for 
comparing the current state of a system or object with a desired or
intended state and then taking the necessary actions 
to transition from the current state to the desired state.
The reconciler is typically used to ensure that the system remains 
in a consistent and valid state 
and that any discrepancies between the current state and 
the desired state are resolved correctly. 
This could involve updating variables, 
executing specific actions, 
or performing other tasks 
to bring the system back into alignment with the desired state.

- `Action` \
An action within a state chart has the capability to interact with 
both the current context and the state of the system. 
By leveraging the context, 
an action can establish communication with external entities like
 APIs, databases, and the file system. 
 Furthermore, the action can manipulate the state, 
 enabling it to read and modify stored information. 
 Should an error occur during the execution of an action, 
 it will interrupt the sequence of subsequent actions within the state, 
 prompting the system to search for an error transition
 to address the issue.

- `Guard` \
A guard within a state chart has the sole ability to read the current state
and verify whether certain conditions are met. 
Typically, these conditions are predicated on the state itself, 
and the guard evaluates them to yield a boolean result. 
When a guard is incorporated into a transition, 
it serves to determine whether the transition should occur, 
thus dictating which state becomes the subsequent one in the sequence.

- `Transition` \
In the context of state charts or state machines, 
a "transition" refers to the process of moving from one state to another. 
A transition is usually triggered by an event or a set of conditions that
are evaluated by a guard. 
When a transition is triggered, 
the system changes its state accordingly. 
Transitions can also have associated actions, 
which are executed when the transition occurs. 
Transitions are an essential concept in state charts 
as they define how the system responds to different events and conditions.

- `(Extended) State` \
The extended state, also known as the context or state data, 
encompasses all the information that is either read or modified
 by actions within a state chart, 
 and also read by guards. 
 Essentially, it serves as the primary point of interaction 
 between various components of the state chart, 
 enabling actions and guards to access, manipulate, 
 and evaluate the current state of the system. 
 The extended state plays a pivotal role in maintaining the integrity and
 functionality of the state chart, 
 ensuring that actions are executed, transitions are evaluated, 
 and guards make informed decisions based on the system's state.

- `Context` \
The "context" serves as a container for storing and managing 
essential information necessary for the interaction 
between a system and its external environment. 
This information typically includes initialized clients, 
which are used to facilitate communication and interaction 
with the outside world, such as APIs, databases, network services, 
or other external systems. 
The context acts as a central repository for these initialized clients,
making them readily accessible to various components within the system, 
such as actions and guards in a state chart. 
By encapsulating this crucial information, 
the context enables efficient and effective communication and
coordination between the system and its external dependencies.
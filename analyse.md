# Endorsement and threat models

## Endorsement on Generating Proof
### Flow
Generate a proof of a campaign c for a user u is as follows:
    - ask authorized verifiers to generate commitment
        + generate random value
        + calculate commitment based on generated value and stored s value for c
        + return commitment and random value

    - calculate total commitment

### Which organizations must participate to validate the proof?
1/ All organizations.
- Safetest
- Slowest
- Unrelated advertiers and businesses must validate the proof
- Use default endorsement policy
- Inconsistency of verifiers' random values.

2/ Only peers of c's advertier and business
- Use key-level endorsement policy
- Can improve performance
- Inconsistency of verifiers' random values.

3/ At least one orgnization
- Fastest
-
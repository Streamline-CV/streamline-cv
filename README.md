# Streamline CV CLI

This CLI tool is designed to assist in managing your CV routine efficiently.

## AI Reviewer

The AI Reviewer leverages artificial intelligence to enhance the content of your CV. It operates by comparing the
current version of `config.yaml` with the main version. For every detected change, the AI is consulted to improve the
text quality. Upon gathering AI feedback, the improvements are documented
in [RDF format](https://github.com/reviewdog/reviewdog/blob/master/proto/rdf/README.md),
enabling [reviewdog](https://github.com/reviewdog/reviewdog) to present these refinements as suggestions on GitHub.

**Note**: The AI Reviewer utilizes the GPT-4 model from OpenAI. Usage incurs a fee, so it's important to use this
feature judiciously. Typically, the costs are minimal, given the brevity expected of a well-crafted CV.

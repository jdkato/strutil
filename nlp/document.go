package nlp

// A DocOpt represents a setting that changes the document creation process.
//
// For example, it might disable named-entity extraction:
//
//	doc := prose.NewDocument("...", prose.WithExtraction(false))
type DocOpt func(doc *Document, opts *DocOpts)

// DocOpts controls the Document creation process:
type DocOpts struct {
	Segment   bool      // If true, include segmentation
	Tag       bool      // If true, include POS tagging
	Tokenizer Tokenizer // If true, include tokenization
}

// UsingTokenizer specifies the Tokenizer to use.
func UsingTokenizer(include Tokenizer) DocOpt {
	return func(doc *Document, opts *DocOpts) {
		// Tagging and entity extraction both require tokenization.
		opts.Tokenizer = include
	}
}

// WithTokenization can enable (the default) or disable tokenization.
// Deprecated: use UsingTokenizer instead.
func WithTokenization(include bool) DocOpt {
	return func(doc *Document, opts *DocOpts) {
		if !include {
			opts.Tokenizer = nil
		}
	}
}

// WithTagging can enable (the default) or disable POS tagging.
func WithTagging(include bool) DocOpt {
	return func(doc *Document, opts *DocOpts) {
		opts.Tag = include
	}
}

// WithSegmentation can enable (the default) or disable sentence segmentation.
func WithSegmentation(include bool) DocOpt {
	return func(doc *Document, opts *DocOpts) {
		opts.Segment = include
	}
}

// UsingModel can enable (the default) or disable named-entity extraction.
func UsingModel(model *Model) DocOpt {
	return func(doc *Document, opts *DocOpts) {
		doc.Model = model
	}
}

// A Document represents a parsed body of text.
type Document struct {
	Model *Model
	Text  string

	sentences []Sentence
	tokens    []*Token
}

// Tokens returns `doc`'s tokens.
func (doc *Document) Tokens() []Token {
	tokens := make([]Token, 0, len(doc.tokens))
	for _, tok := range doc.tokens {
		tokens = append(tokens, *tok)
	}
	return tokens
}

// Sentences returns `doc`'s sentences.
func (doc *Document) Sentences() []Sentence {
	return doc.sentences
}

var defaultOpts = DocOpts{
	Tokenizer: NewIterTokenizer(),
	Segment:   true,
	Tag:       true,
}

// NewDocument creates a Document according to the user-specified options.
//
// For example,
//
//	doc := prose.NewDocument("...")
func NewDocument(text string, opts ...DocOpt) (*Document, error) {
	var pipeError error

	doc := Document{Text: text}
	base := defaultOpts
	for _, applyOpt := range opts {
		applyOpt(&doc, &base)
	}

	if doc.Model == nil {
		doc.Model = defaultModel(base.Tag)
	}

	if base.Segment {
		segmenter := newPunktSentenceTokenizer()
		doc.sentences = segmenter.segment(text)
	}
	if base.Tokenizer != nil {
		doc.tokens = append(doc.tokens, base.Tokenizer.Tokenize(text)...)
	}
	if base.Tag {
		doc.tokens = doc.Model.tagger.tag(doc.tokens)
	}

	return &doc, pipeError
}

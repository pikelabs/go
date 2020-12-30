package rpm

type getTagFn func(h *Header, tag HeaderTag) (*HeaderIndexEntry, []byte, error)

const (
	TagFilenames = 5000
)

var extTagFunc = map[HeaderTag]getTagFn{
	TagFilenames: getTagFilenames,
}

func getTagFilenames(h *Header, tag HeaderTag) (*HeaderIndexEntry, []byte, error) {
	return getTag(h, TagBaseNames)
}

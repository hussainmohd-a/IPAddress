package ipaddr

type SegmentsIterator interface {
	iteratorBase
	Next() []*AddressDivision
}

type singleSegmentsIterator struct {
	original []*AddressDivision
}

func (it *singleSegmentsIterator) HasNext() bool {
	return it.original != nil
}

func (it *singleSegmentsIterator) Next() (res []*AddressDivision) {
	if it.HasNext() {
		res = it.original
		it.original = nil
	}
	return
}

type multiSegmentsIterator struct {
	done       bool
	variations []SegmentIterator
	nextSet    []*AddressDivision

	segIteratorProducer,
	hostSegIteratorProducer func(int) SegmentIterator

	networkSegmentIndex,
	hostSegmentIndex int

	excludeFunc func([]*AddressDivision) bool
}

func (it *multiSegmentsIterator) HasNext() bool {
	return !it.done
}

func (it *multiSegmentsIterator) updateVariations(start int) {
	i := start
	nextSet := it.nextSet
	variations := it.variations
	segIteratorProducer := it.segIteratorProducer
	for ; i < it.hostSegmentIndex; i++ {
		variations[i] = segIteratorProducer(i)
		nextSet[i] = variations[i].Next().ToAddressDivision()
	}
	if i == it.networkSegmentIndex {
		variations[i] = it.hostSegIteratorProducer(i)
		nextSet[i] = variations[i].Next().ToAddressDivision()
	}
}

func (it *multiSegmentsIterator) init() {
	it.updateVariations(0)
	nextSet := it.nextSet
	variations := it.variations
	divCount := len(variations)
	hostSegIteratorProducer := it.hostSegIteratorProducer
	// for regular iterators (not prefix block), networkSegmentIndex is last segment (count - 1)
	for i := it.networkSegmentIndex + 1; i < divCount; i++ {
		variations[i] = hostSegIteratorProducer(i)
		nextSet[i] = variations[i].Next().ToAddressDivision()
	}
	excludeFunc := it.excludeFunc
	if excludeFunc != nil && excludeFunc(it.nextSet) {
		it.increment()
	}
}

func (it *multiSegmentsIterator) Next() (res []*AddressDivision) {
	if it.HasNext() {
		res = it.increment()
	}
	return
}

func (it *multiSegmentsIterator) increment() (res []*AddressDivision) {
	var previousSegs []*AddressDivision
	// the current set of segments already holds the next iteration,
	// this searches for the set of segments to follow.
	variations := it.variations
	nextSet := it.nextSet
	for j := it.networkSegmentIndex; j >= 0; j-- { //for regular iterators (not prefix block), networkSegmentIndex is last segment (count - 1)
		for variations[j].HasNext() {
			if previousSegs == nil {
				previousSegs = cloneDivs(nextSet)
			}
			nextSet[j] = variations[j].Next().ToAddressDivision()
			it.updateVariations(j + 1)
			excludeFunc := it.excludeFunc
			if excludeFunc != nil && excludeFunc(nextSet) {
				// try again, starting over
				j = it.networkSegmentIndex
			} else {
				return previousSegs
			}
		}
	}
	it.done = true
	if previousSegs == nil {
		// never found set of candidate segments
		return nextSet
	}
	// found a candidate to follow, but was rejected.
	// nextSet has that rejected candidate,
	// so we must return the set that was created prior to that.
	return previousSegs
}

// this iterator function used by addresses and segment arrays, for iterators that are not prefix or prefix block iterators
func allSegmentsIterator(
	divCount int,
	segSupplier func() []*AddressDivision, // only useful for a segment iterator.  Address/section iterators use address/section for single valued iterator.
	segIteratorProducer func(int) SegmentIterator,
	excludeFunc func([]*AddressDivision) bool /* can be nil */) SegmentsIterator {
	return segmentsIterator(divCount, segSupplier, segIteratorProducer, excludeFunc, divCount-1, divCount, nil)
}

// used to produce regular iterators with or without zero-host values, and prefix block iterators
func segmentsIterator(
	divCount int,
	segSupplier func() []*AddressDivision,
	segIteratorProducer func(int) SegmentIterator, // unused at this time, since we do not have a public segments iterator
	excludeFunc func([]*AddressDivision) bool, // can be nil
	networkSegmentIndex,
	hostSegmentIndex int,
	hostSegIteratorProducer func(int) SegmentIterator) SegmentsIterator { // returns Iterator<S[]>
	if segSupplier != nil {
		return &singleSegmentsIterator{segSupplier()}
	}
	iterator := &multiSegmentsIterator{
		variations:              make([]SegmentIterator, divCount),
		nextSet:                 make([]*AddressDivision, divCount),
		segIteratorProducer:     segIteratorProducer,
		hostSegIteratorProducer: hostSegIteratorProducer,
		networkSegmentIndex:     networkSegmentIndex,
		hostSegmentIndex:        hostSegmentIndex,
		excludeFunc:             excludeFunc,
	}
	iterator.init()
	return iterator
}

// this iterator function used by sequential ranges
func rangeSegmentsIterator(
	divCount int,
	segIteratorProducer func(int) SegmentIterator,
	networkSegmentIndex,
	hostSegmentIndex int,
	prefixedSegIteratorProducer func(int) SegmentIterator) SegmentsIterator {
	return segmentsIterator(
		divCount,
		nil,
		segIteratorProducer,
		nil,
		networkSegmentIndex,
		hostSegmentIndex,
		prefixedSegIteratorProducer)
}

type SectionIterator interface {
	HasNext() bool
	Next() *AddressSection
}

type singleSectionIterator struct {
	original *AddressSection
}

func (it *singleSectionIterator) HasNext() bool {
	return it.original != nil
}

func (it *singleSectionIterator) Next() (res *AddressSection) {
	if it.HasNext() {
		res = it.original
		it.original = nil
	}
	return
}

type multiSectionIterator struct {
	original        *AddressSection
	iterator        SegmentsIterator
	valsAreMultiple bool
}

func (it *multiSectionIterator) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *multiSectionIterator) Next() (res *AddressSection) {
	if it.HasNext() {
		segs := it.iterator.Next()
		original := it.original
		res = createSection(segs, original.prefixLength, original.addrType)
		res.isMultiple = it.valsAreMultiple
	}
	return
}

func sectIterator(
	useOriginal bool,
	original *AddressSection,
	valsAreMultiple bool,
	iterator SegmentsIterator,
) SectionIterator {
	if useOriginal {
		return &singleSectionIterator{original: original}
	}
	return &multiSectionIterator{
		original:        original,
		iterator:        iterator,
		valsAreMultiple: valsAreMultiple,
	}
}

type prefixSectionIterator struct {
	original   *AddressSection
	iterator   SegmentsIterator
	isNotFirst bool
}

func (it *prefixSectionIterator) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *prefixSectionIterator) Next() (res *AddressSection) {
	if it.HasNext() {
		segs := it.iterator.Next()
		original := it.original
		res = createSection(segs, original.prefixLength, original.addrType)
		if !it.isNotFirst {
			res.initMultiple() // sets isMultiple
			it.isNotFirst = true
		} else if !it.HasNext() {
			res.initMultiple() // sets isMultiple
		} else {
			res.isMultiple = true
		}
	}
	return
}

func prefixSectIterator(
	useOriginal bool,
	original *AddressSection,
	iterator SegmentsIterator,
) SectionIterator {
	if useOriginal {
		return &singleSectionIterator{original: original}
	}
	return &prefixSectionIterator{
		original: original,
		iterator: iterator,
	}
}

// IPSectionIterator iterates through IP address and subnet sections
type IPSectionIterator interface {
	iteratorBase
	Next() *IPAddressSection
}

type ipSectionIterator struct {
	SectionIterator
}

func (iter ipSectionIterator) Next() *IPAddressSection {
	return iter.SectionIterator.Next().ToIPAddressSection()
}

// IPv4SectionIterator iterates through IPv4 address and subnet sections
type IPv4SectionIterator interface {
	iteratorBase
	Next() *IPv4AddressSection
}

type ipv4SectionIterator struct {
	SectionIterator
}

func (iter ipv4SectionIterator) Next() *IPv4AddressSection {
	return iter.SectionIterator.Next().ToIPv4AddressSection()
}

// IPv6SectionIterator iterates through IPv6 address and subnet sections
type IPv6SectionIterator interface {
	iteratorBase
	Next() *IPv6AddressSection
}

type ipv6SectionIterator struct {
	SectionIterator
}

func (iter ipv6SectionIterator) Next() *IPv6AddressSection {
	return iter.SectionIterator.Next().ToIPv6AddressSection()
}

// MACSectionIterator iterates through MAC address and subnet sections
type MACSectionIterator interface {
	iteratorBase
	Next() *MACAddressSection
}

type macSectionIterator struct {
	SectionIterator
}

func (iter macSectionIterator) Next() *MACAddressSection {
	return iter.SectionIterator.Next().ToMACAddressSection()
}

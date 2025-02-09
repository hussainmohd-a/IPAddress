/*
 * Copyright 2018 Sean C Foley
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *     or at
 *     https://github.com/seancfoley/IPAddress/blob/master/LICENSE
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package inet.ipaddr;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Iterator;
import java.util.NoSuchElementException;
import java.util.function.BiFunction;
import java.util.function.Function;
import java.util.function.IntFunction;
import java.util.function.Predicate;
import java.util.function.Supplier;
import java.util.function.ToLongFunction;
import java.util.function.UnaryOperator;
import java.util.stream.Stream;

import inet.ipaddr.AddressNetwork.AddressSegmentCreator;
import inet.ipaddr.IPAddressSection.IPAddressSeqRangeSpliterator;
import inet.ipaddr.IPAddressSection.IPAddressSeqRangePrefixSpliterator;
import inet.ipaddr.IPAddressSection.SegFunction;
import inet.ipaddr.IPAddressSection.SeqRangeIteratorProvider;
import inet.ipaddr.format.AddressComponentRange;
import inet.ipaddr.format.IPAddressRange;
import inet.ipaddr.format.standard.AddressCreator;
import inet.ipaddr.format.util.AddressComponentRangeSpliterator;
import inet.ipaddr.format.util.AddressComponentSpliterator;
import inet.ipaddr.format.validate.ParsedAddressGrouping;

/**
 * This class can be used to represent an arbitrary range of IP addresses.  
 * <p>
 * Note that the IPAddress and IPAddressString classes allow you to specify a range of values for each segment.
 * That allows you to represent single addresses, any address prefix subnet (eg 1.2.0.0/16 or 1:2:3:4::/64) or any subnet that can be represented with segment ranges (1.2.0-255.* or 1:2:3:4:*), see
 * {@link IPAddressString} for details.
 * <p>
 * IPAddressString and IPAddress cover all potential subnets and addresses that can be represented by a single address string of 4 or less segments for IPv4, and 8 or less segments for IPv6.
 * <p>
 * This class allows the representation of any sequential address range, including those that cannot be represented by IPAddress.
 * <p>
 * String representations include the full address for both the lower and upper bounds of the range.
 *  
 * @custom.core
 * @author sfoley
 *
 */
public abstract class IPAddressSeqRange implements IPAddressRange {
	
	private static final long serialVersionUID = 1L;
	
	protected final IPAddress lower, upper;
	
	private transient BigInteger count;
	private transient int hashCode;

	protected <T extends IPAddress> IPAddressSeqRange(
			T first, 
			T other,
			UnaryOperator<T> getLower,
			UnaryOperator<T> getUpper,
			UnaryOperator<T> prefixLenRemover) {
		boolean f;
		if((f = first.contains(other)) || other.contains(first)) {
			T addr = f ? prefixLenRemover.apply(first) : prefixLenRemover.apply(other);
			lower = getLower.apply(addr);
			upper = getUpper.apply(addr);
		} else {
			T firstLower = getLower.apply(first);
			T otherLower = getLower.apply(other);
			T firstUpper = getUpper.apply(first);
			T otherUpper = getUpper.apply(other);
			T lower = compareLowValues(firstLower, otherLower) > 0 ? otherLower : firstLower;
			T upper = compareLowValues(firstUpper, otherUpper) < 0 ? otherUpper : firstUpper;
			this.lower = prefixLenRemover.apply(lower);
			this.upper = prefixLenRemover.apply(upper);
		}
	}
	
	protected <T extends IPAddress> IPAddressSeqRange(
			T first, 
			T second) {
		lower = first;
		upper = second;
	}

	private static int compareLowValues(IPAddress one, IPAddress two) {
		return IPAddress.compareLowValues(one, two);
	}
	
	@Override
	public BigInteger getCount() {
		BigInteger result = count;
		if(result == null) {
			count = result = getCountImpl();
		}
		return result;
	}
	
	@Override
	public boolean isMultiple() {
		BigInteger count = this.count;
		if(count == null) {
			return !getLower().equals(getUpper());
		}
		return IPAddressRange.super.isMultiple();
	}
	
	/**
	 * 
	 * @param other the range to compare, which does not need to range across the same address space
	 * @return whether this range spans more addresses than the provided range.
	 */
	public boolean isMore(IPAddressSeqRange other) {
		return getCount().compareTo(other.getCount()) > 0;
	}
	
	protected BigInteger getCountImpl() {
		return IPAddressRange.super.getCount();
	}
	
	@Override
	public abstract Iterable<? extends IPAddress> getIterable();
	
	protected static int getNetworkSegmentIndex(int networkPrefixLength, int bytesPerSegment, int bitsPerSegment) {
		return ParsedAddressGrouping.getNetworkSegmentIndex(networkPrefixLength, bytesPerSegment, bitsPerSegment);
	}
	
	protected static int getHostSegmentIndex(int networkPrefixLength, int bytesPerSegment, int bitsPerSegment) {
		return ParsedAddressGrouping.getHostSegmentIndex(networkPrefixLength, bytesPerSegment, bitsPerSegment);
	}
	
	/**
	 * Iterates through the range of prefix blocks in this range instance using the given prefix length.
	 * 
	 * @param prefLength
	 * @return
	 */
	@Override
	public abstract Iterator<? extends IPAddress> prefixBlockIterator(int prefLength);

	@Override
	public abstract AddressComponentRangeSpliterator<? extends IPAddressSeqRange, ? extends IPAddress> prefixBlockSpliterator(int prefLength);

	@Override
	public abstract Stream<? extends IPAddress> prefixBlockStream(int prefLength);

	protected static interface IPAddressSeqRangeSplitterSink<S, T>{
		void setSplitValues(S left, S right);
		
		S getAddressItem();
	};
	
	@FunctionalInterface
	protected static interface IPAddressSeqRangeIteratorProvider<S, T> extends SeqRangeIteratorProvider<S,T>{}
	
	protected static <S extends AddressComponentRange, T> AddressComponentRangeSpliterator<S, T> createSpliterator(
			S forIteration,
			Predicate<IPAddressSeqRangeSplitterSink<S, T>> splitter,
			IPAddressSeqRangeIteratorProvider<S, T> iteratorProvider,
			ToLongFunction<S> longSizer) {
		return new IPAddressSeqRangeSpliterator<S, T>(forIteration, splitter, iteratorProvider, longSizer);
	}
	
	protected static <S extends AddressComponentRange, T> AddressComponentRangeSpliterator<S, T> createSpliterator(
			S forIteration,
			Predicate<IPAddressSeqRangeSplitterSink<S, T>> splitter,
			IPAddressSeqRangeIteratorProvider<S, T> iteratorProvider,
			Function<S, BigInteger> sizer,
			Predicate<S> downSizer,
			ToLongFunction<S> longSizer) {
		return new IPAddressSeqRangeSpliterator<S, T>(forIteration, splitter, iteratorProvider, sizer, downSizer, longSizer);
	}

	protected static <S extends AddressComponentRange> AddressComponentSpliterator<S> createPrefixSpliterator(
			S forIteration,
			Predicate<IPAddressSeqRangeSplitterSink<S, S>> splitter,
			IPAddressSeqRangeIteratorProvider<S, S> iteratorProvider,
			ToLongFunction<S> longSizer) {
		return new IPAddressSeqRangePrefixSpliterator<S>(forIteration, splitter, iteratorProvider, longSizer);
	}
	
	protected static <S extends AddressComponentRange> AddressComponentSpliterator<S> createPrefixSpliterator(
			S forIteration,
			Predicate<IPAddressSeqRangeSplitterSink<S, S>> splitter,
			IPAddressSeqRangeIteratorProvider<S, S> iteratorProvider,
			Function<S, BigInteger> sizer,
			Predicate<S> downSizer,
			ToLongFunction<S> longSizer) {
		return new IPAddressSeqRangePrefixSpliterator<S>(forIteration, splitter, iteratorProvider, sizer, downSizer, longSizer);
	}
	
	protected static <R,A extends IPAddress> Iterator<R> rangedIterator(Iterator<A> iter) {
		return new Iterator<R>() {
			@Override
			public boolean hasNext() {
				return iter.hasNext();
			}

			@SuppressWarnings("unchecked")
			@Override
			public R next() {
				return (R) iter.next().toSequentialRange();
			}
		};
	}

	/**
	 * Iterates through the range of prefixes in this range instance using the given prefix length.
	 * <p>
	 * Since a range between two arbitrary addresses cannot always be represented with a single IPAddress instance,
	 * the returned iterator iterates through IPAddressRange instances.
	 * <p>
	 * For instance, if iterating from 1.2.3.4 to 1.2.4.5 with prefix 8, the range shares the same prefix 1,
	 * but the range cannot be represented by the address 1.2.3-4.4-5 which does not include 1.2.3.255 or 1.2.4.0 both of which are in the original range.
	 * Nor can the range be represented by 1.2.3-4.0-255 which includes 1.2.4.6 and 1.2.3.3, both of which were not in the original range.
	 * An IPAddressSeqRange is thus required to represent that prefixed range.
	 * 
	 * @param prefixLength
	 * @return
	 */
	@Override
	public Iterator<? extends IPAddressSeqRange> prefixIterator(int prefixLength) {
		if(prefixLength < 0) {
			throw new PrefixLenException(prefixLength);
		}
		if(!isMultiple()) {
			return new Iterator<IPAddressSeqRange>() {
				IPAddressSeqRange orig = IPAddressSeqRange.this;

				@Override
				public boolean hasNext() {
					return orig != null;
				}

			    @Override
				public IPAddressSeqRange next() {
			    	if(orig == null) {
			    		throw new NoSuchElementException();
			    	}
			    	IPAddressSeqRange result = orig;
			    	orig = null;
			    	return result;
			    }
			
			    @Override
				public void remove() {
			    	throw new UnsupportedOperationException();
			    }
			};
		}
		return new Iterator<IPAddressSeqRange>() {
			Iterator<? extends IPAddress> prefixBlockIterator = prefixBlockIterator(prefixLength);
			private boolean first = true;

			@Override
			public boolean hasNext() {
				return prefixBlockIterator.hasNext();
			}

		    @Override
			public IPAddressSeqRange next() {
		    	IPAddress next = prefixBlockIterator.next();
		    	if(first) {
		    		first = false;
		    		// next is a prefix block
		    		IPAddress lower = getLower();
		    		if(hasNext()) {
			    		if(!lower.includesZeroHost(prefixLength)) {
			    			return create(lower, next.getUpper());
			    		}
		    		} else {
		    			IPAddress upper = getUpper();
		    			if(!lower.includesZeroHost(prefixLength) || !upper.includesMaxHost(prefixLength)) {
		    				return create(lower, upper);
		    			}
		    		}
		    	} else if(!hasNext()) {
		    		IPAddress upper = getUpper();
		    		if(!upper.includesMaxHost(prefixLength)) {
		    			return create(next.getLower(), upper);
		    		}
		    	}
		    	return next.toSequentialRange();
		    }
		
		    @Override
			public void remove() {
		    	throw new UnsupportedOperationException();
		    }
		};
	}

	@Override
	public abstract AddressComponentSpliterator<? extends IPAddressSeqRange> prefixSpliterator(int prefLength);
	
	@Override
	public abstract Stream<? extends IPAddressSeqRange> prefixStream(int prefLength);

	@FunctionalInterface
	protected interface SegValueComparator<T> {
	    boolean apply(T segmentSeries1, T segmentSeries2, int index);
	}

	/**
	 * Splits a sequential range into two
	 * <p>
	 * Returns false if it cannot be done
	 * 
	 * @param beingSplit
	 * @param transformer
	 * @param segmentCreator
	 * @param originalSegments
	 * @param networkSegmentIndex if this index matches hostSegmentIndex, splitting will attempt to split the network part of this segment
	 * @param hostSegmentIndex splitting will work with the segments prior to this one
	 * @param prefixLength
	 * @return
	 */
	protected static <I extends IPAddressSeqRange, T extends IPAddressRange, S extends AddressSegment> boolean split(
			IPAddressSeqRangeSplitterSink<I, T> sink,
			BiFunction<S[], S[], I> transformer,
			AddressSegmentCreator<S> segmentCreator,
			S originalSegmentsLower[],
			S originalSegmentsUpper[],
			int networkSegmentIndex, //for regular iterators (not prefix block), networkSegmentIndex is last segment (count - 1) - it is only instrumental with prefix iterators
			int hostSegmentIndex, // for regular iterators hostSegmentIndex is past last segment (count) - it is only instrumental with prefix iterators
			Integer prefixLength) {
		int i = 0;
		S lowerSeg, upperSeg;
		lowerSeg = upperSeg = null;
		boolean isSplit = false;
		for(; i < hostSegmentIndex; i++) {
			S segLower = originalSegmentsLower[i];
			S segUpper = originalSegmentsUpper[i];
			int lower = segLower.getSegmentValue();
			int upper = segUpper.getSegmentValue();
			// if multiple, split into two
			if(lower != upper) {
				isSplit = true;
				int size = upper - lower;
				int mid = lower + (size >>> 1);
				lowerSeg = segmentCreator.createSegment(mid);
				upperSeg = segmentCreator.createSegment(mid + 1);
				break;
			}
		}
		if(i == networkSegmentIndex && !isSplit) {
			// prefix or prefix block iterators: no need to differentiate, handle both as prefix, iteration will handle the rest
			S segLower = originalSegmentsLower[i];
			S segUpper = originalSegmentsUpper[i];
			int segBitCount = segLower.getBitCount();
			Integer pref = IPAddressSection.getSegmentPrefixLength(segBitCount, prefixLength, i);
			int shiftAdjustment = segBitCount - pref;
			int lower = segLower.getSegmentValue();
			int upper = segUpper.getSegmentValue();
			lower >>>= shiftAdjustment;
			upper >>>= shiftAdjustment;
			if(lower != upper) {
				isSplit = true;
				int size = upper - lower;
				int mid = lower + (size >>> 1);
				int next = mid + 1;
				mid = (mid << shiftAdjustment) | ~(~0 << shiftAdjustment);
				next <<= shiftAdjustment;
				lowerSeg = segmentCreator.createSegment(mid);
				upperSeg = segmentCreator.createSegment(next);
			}
		}
		if(isSplit) {
			int len = originalSegmentsLower.length;
			S lowerUpperSegs[] = segmentCreator.createSegmentArray(len);
			S upperLowerSegs[] = segmentCreator.createSegmentArray(len);
			System.arraycopy(originalSegmentsLower, 0, lowerUpperSegs, 0, i);
			System.arraycopy(originalSegmentsLower, 0, upperLowerSegs, 0, i);
			int j = i + 1;
			lowerUpperSegs[i] = lowerSeg;
			upperLowerSegs[i] = upperSeg;
			Arrays.fill(lowerUpperSegs, j, lowerUpperSegs.length, segmentCreator.createSegment(lowerSeg.getMaxSegmentValue()));
			Arrays.fill(upperLowerSegs, j, upperLowerSegs.length, segmentCreator.createSegment(0));
			sink.setSplitValues(transformer.apply(originalSegmentsLower, lowerUpperSegs), transformer.apply(upperLowerSegs, originalSegmentsUpper));
		}
		return isSplit;
	}

	@Override
	public abstract Iterator<? extends IPAddress> iterator();

	@Override
	public abstract AddressComponentRangeSpliterator<? extends IPAddressSeqRange, ? extends IPAddress> spliterator();
	
	@Override
	public abstract Stream<? extends IPAddress> stream();

	/*
	 * This iterator is used for the case where the range is non-multiple
	 */
	protected static <T extends Address, S extends AddressSegment> Iterator<T> iterator(T original, AddressCreator<T, ?, ?, S> creator) {
		return IPAddressSection.iterator(original, creator, null);
	}
	
	/*
	 This iterator is (not surprisingly) 2 to 3 times faster (based on measurements I've done) than an iterator that uses the increment method like:
	 
	 return iterator(a -> a.increment(1));
	 
	 protected Iterator<T> iterator(UnaryOperator<T> incrementor) {
	 	return new Iterator<T>() {
			BigInteger count = getCount();
			T current = lower;
					
			@Override
			public boolean hasNext() {
				return !count.equals(BigInteger.ZERO);
			}

			@Override
			public T next() {
				if(hasNext()) {
					T result = current;
					current = incrementor.apply(current);
					count = count.subtract(BigInteger.ONE);
					return result;
				}
				throw new NoSuchElementException();
			}
		};
	 }
	 */
	protected static <T extends IPAddress, S extends IPAddressSegment> Iterator<T> iterator(
			T lower,
			T upper,
			AddressCreator<T, ?, ?, S> creator,
			SegFunction<T, S> segProducer,
			SegFunction<S, Iterator<S>> segmentIteratorProducer,
			SegValueComparator<T> segValueComparator,
			int networkSegmentIndex,
			int hostSegmentIndex,
			SegFunction<S, Iterator<S>> prefixedSegIteratorProducer) {
		int divCount = lower.getDivisionCount();
		
		// at any given point in time, this list provides an iterator for the segment at each index
		ArrayList<Supplier<Iterator<S>>> segIteratorProducerList = new ArrayList<Supplier<Iterator<S>>>(divCount);
		
		// at any given point in time, finalValue[i] is true if and only if we have reached the very last value for segment i - 1
		// when that happens, the next iterator for the segment at index i will be the last
		boolean finalValue[] = new boolean[divCount + 1];
		
		// here is how the segment iterators will work:
		// the low and high values at each segment are low, high
		// the maximum possible values for any segment are min, max
		// we first find the first k >= 0 such that low != high for the segment at index k
		
		//	the initial set of iterators at each index are as follows:
		//    for i < k finalValue[i] is set to true right away.
		//		we create an iterator from seg = new Seg(low)
		//    for i == k we create a wrapped iterator from Seg(low, high), wrapper will set finalValue[i] once we reach the final value of the iterator
		//    for i > k we create an iterator from Seg(low, max)
		// 
		// after the initial iterator has been supplied, any further iterator supplied for the same segment is as follows:
		//    for i <= k, there was only one iterator, there will be no further iterator
		//    for i > k,
		//	  	if i == 0 or of if flagged[i - 1] is true, we create a wrapped iterator from Seg(low, high), wrapper will set finalValue[i] once we reach the final value of the iterator
		//      otherwise we create an iterator from Seg(min, max)
		//
		// By following these rules, we iterate through all possible addresses	
		boolean notDiffering = true;
		finalValue[0] = true;
		S allSegShared = null;
		for(int i = 0; i < divCount; i++) {
			SegFunction<S, Iterator<S>> segIteratorProducer;
			if(prefixedSegIteratorProducer != null && i >= networkSegmentIndex) {
				segIteratorProducer = prefixedSegIteratorProducer;
			} else {
				segIteratorProducer = segmentIteratorProducer;
			}
			S lowerSeg = segProducer.apply(lower, i);
			int indexi = i;
			if(notDiffering) {
				notDiffering = segValueComparator.apply(lower, upper, i);
				if(notDiffering) {
					// there is only one iterator and it produces only one value
					finalValue[i + 1] = true;
					Iterator<S> iterator = segIteratorProducer.apply(lowerSeg, i);
					segIteratorProducerList.add(() -> iterator);
				} else {
					// in the first differing segment the only iterator will go from segment value of lower address to segment value of upper address
					Iterator<S> iterator = segIteratorProducer.apply(
							creator.createSegment(lowerSeg.getSegmentValue(), upper.getSegment(i).getSegmentValue(), null), i);
					
					// the wrapper iterator detects when the iterator has reached its final value
					Iterator<S> wrappedFinalIterator = new Iterator<S>() {
						@Override
						public boolean hasNext() {
							return iterator.hasNext();
						}

						@Override
						public S next() {
							S next = iterator.next();
							if(!iterator.hasNext()) {
								finalValue[indexi + 1] = true;
							}
							return next;
						}
					};
					segIteratorProducerList.add(() -> wrappedFinalIterator);
				}
			} else {
				// in the second and all following differing segments, rather than go from segment value of lower address to segment value of upper address
				// we go from segment value of lower address to the max seg value the first time through
				// then we go from the min value of the seg to the max seg value each time until the final time,
				// the final time we go from the min value to the segment value of upper address
				// we know it is the final time through when the previous iterator has reached its final value, which we track
				
				// the first iterator goes from the segment value of lower address to the max value of the segment
				Iterator<S> firstIterator = segIteratorProducer.apply(creator.createSegment(lowerSeg.getSegmentValue(), lower.getMaxSegmentValue(), null), i);
				
				// the final iterator goes from 0 to the segment value of our upper address
				Iterator<S> finalIterator = segIteratorProducer.apply(creator.createSegment(0, upper.getSegment(i).getSegmentValue(), null), i);
				
				// the wrapper iterator detects when the final iterator has reached its final value
				Iterator<S> wrappedFinalIterator = new Iterator<S>() {
					@Override
					public boolean hasNext() {
						return finalIterator.hasNext();
					}

					@Override
					public S next() {
						S next = finalIterator.next();
						if(!finalIterator.hasNext()) {
							finalValue[indexi + 1] = true;
						}
						return next;
					}
				};
				if(allSegShared == null) {
					allSegShared = creator.createSegment(0, lower.getMaxSegmentValue(), null);
				}
				// all iterators after the first iterator and before the final iterator go from 0 the max segment value,
				// and there will be many such iterators
				S allSeg = allSegShared;
				Supplier<Iterator<S>> finalIteratorProducer = () -> finalValue[indexi] ?  wrappedFinalIterator : segIteratorProducer.apply(allSeg, indexi);
				segIteratorProducerList.add(() -> {
					//the first time through, we replace the iterator producer so the first iterator used only once
					segIteratorProducerList.set(indexi, finalIteratorProducer);
					return firstIterator;
				});
			}
		}
		IntFunction<Iterator<S>> iteratorProducer = iteratorIndex -> segIteratorProducerList.get(iteratorIndex).get();
		return IPAddressSection.iterator(null, creator,
				IPAddressSection.iterator(
						lower.getSegmentCount(),
						creator,
						iteratorProducer, 
						networkSegmentIndex,
						hostSegmentIndex,
						iteratorProducer)
			);
	}
	
	@Override
	public IPAddress getLower() {
		return lower;
	}
	
	@Override
	public IPAddress getUpper() {
		return upper;
	}
	
	public String toNormalizedString(String separator) {
		Function<IPAddress, String> stringer = IPAddress::toNormalizedString;
		return toString(stringer, separator, stringer);
	}
	
	@Override
	public String toNormalizedString() {
		return toNormalizedString(" -> ");
	}
	
	public String toCanonicalString(String separator) {
		Function<IPAddress, String> stringer = IPAddress::toCanonicalString;
		return toString(stringer, separator, stringer);
	}
	
	@Override
	public String toCanonicalString() {
		return toCanonicalString(" -> ");
	}
	
	public String toString(Function<? super IPAddress, String> lowerStringer, String separator, Function<? super IPAddress, String> upperStringer) {
		return lowerStringer.apply(getLower()) + separator + upperStringer.apply(getUpper());
	}
	
	@Override
	public String toString() {
		return toCanonicalString();
	}

	@Override
	public abstract IPAddress coverWithPrefixBlock();
	
	@Override
	public abstract IPAddress[] spanWithPrefixBlocks();

	@Override
	public abstract IPAddress[] spanWithSequentialBlocks();
	
	/**
	 * Joins the given ranges into the fewest number of ranges.
	 * The returned array will be sorted by ascending lowest range value. 
	 * 
	 * @param ranges
	 * @return
	 */
	public static IPAddressSeqRange[] join(IPAddressSeqRange... ranges) {
		ranges = ranges.clone();
		// null entries are automatic joins
		int joinedCount = 0;
		for(int i = 0, j = ranges.length - 1; i <= j; i++) {
			if(ranges[i] == null) {
				joinedCount++;
				while(ranges[j] == null && j > i) {
					j--;
					joinedCount++;
				}
				if(j > i) {
					ranges[i] = ranges[j];
					ranges[j] = null;
					j--;
				}
			}
		}
		int len = ranges.length - joinedCount;
		Arrays.sort(ranges, 0, len, Address.ADDRESS_LOW_VALUE_COMPARATOR);
		for(int i = 0; i < len; i++) {
			IPAddressSeqRange range = ranges[i];
			if(range == null) {
				continue;
			}
			IPAddress currentLower = range.getLower();
			IPAddress currentUpper = range.getUpper();
			boolean didJoin = false;
			for(int j = i + 1; j < ranges.length; j++) {
				IPAddressSeqRange range2 = ranges[j];
				if(range2 == null) {
					continue;
				}
				IPAddress nextLower = range2.getLower();
				if(compareLowValues(currentUpper, nextLower) >= 0
						|| currentUpper.increment(1).equals(nextLower)) {
					//join them
					IPAddress nextUpper = range2.getUpper();
					if(compareLowValues(currentUpper, nextUpper) < 0) {
						currentUpper = nextUpper;
					}
					ranges[j] = null;
					didJoin = true;
					joinedCount++;
				} else break;
			}
			if(didJoin) {
				ranges[i] = range.create(currentLower, currentUpper);
			}
		}
		if(joinedCount == 0) {
			return ranges;
		}
		IPAddressSeqRange joined[] = new IPAddressSeqRange[ranges.length - joinedCount];
		for(int i = 0, j = 0; i < ranges.length; i++) {
			IPAddressSeqRange range = ranges[i];
			if(range == null) {
				continue;
			}
			joined[j++] = range;
			if(j >= joined.length) {
				break;
			}
		}
		return joined;
	}
	
	public boolean overlaps(IPAddressSeqRange other) {
		return compareLowValues(other.getLower(), getUpper()) <= 0 && compareLowValues(other.getUpper(), getLower()) >= 0;
	}

	// we choose to not make this public
	// it is simply wrong to do an instanceof in the IPAddressRange interface, it assumes you know all implementors, 
	// it will not work if/when someone adds a new implementation.
	// If you do not know how an IPAddressRange is implemented, can you do the contains?  
	// Yes, but only by iterating, which is damn ugly for large ranges.
	// Now that we have toSequentialRange() in IPAddressRange, it is easy to do this for sequential subnets.
	// And for non-sequential, there is no simple way of doing it, 
	// in IPAddress you need to either go through the segments, or you need to go through the sequential blocks,
	// and there is no general way to do it for any implementation of IPAddressRange.
	private boolean containsRange(IPAddressRange other) {
		return compareLowValues(other.getLower(), getLower()) >= 0 && compareLowValues(other.getUpper(), getUpper()) <= 0;
	}
	
	@Override
	public boolean contains(IPAddress other) {
		return containsRange(other);
	}
	
	@Override
	public boolean contains(IPAddressSeqRange other) {
		return containsRange(other);
	}
	
	@Override
	public boolean isSequential() {
		return true;
	}
	
	@Override
	public int hashCode() {
		int res = hashCode;
		if(res == 0) {
			res = 31 * getLower().hashCode() + getUpper().hashCode();
			hashCode = res;
		}
		return res;
	}
	
	@Override
	public boolean equals(Object o) {
		if(o instanceof IPAddressSeqRange) {
			IPAddressSeqRange otherRange = (IPAddressSeqRange) o;
				return getLower().equals(otherRange.getLower()) && getUpper().equals(otherRange.getUpper());
			}
			return false;
	}
	
	/**
	 * Returns the intersection of this range with the given range, a range which includes those addresses in both this and the given range.
	 * @param other
	 * @return
	 */
	public IPAddressSeqRange intersect(IPAddressSeqRange other) {
		IPAddress otherLower = other.getLower();
		IPAddress otherUpper = other.getUpper();
		IPAddress lower = getLower();
		IPAddress upper = getUpper();
		if(compareLowValues(lower, otherLower) <= 0) {
			if(compareLowValues(upper, otherUpper) >= 0) {
				return other;
			} else if(compareLowValues(upper, otherLower) < 0) {
				return null;
			}
			return create(otherLower, upper);
		} else if(compareLowValues(otherUpper, upper) >= 0) {
			return this;
		} else if(compareLowValues(otherUpper, lower) < 0) {
			return null;
		}
		return create(lower, otherUpper);
	}
	
	/**
	 * If this range overlaps with the given range,
	 * or if the highest value of the lower range is one below the lowest value of the higher range,
	 * then the two are joined into a new larger range that is returned.
	 * <p>
	 * Otherwise null is returned.
	 * 
	 * @param other
	 * @return
	 */
	public IPAddressSeqRange join(IPAddressSeqRange other) {
		IPAddress otherLower = other.getLower();
		IPAddress otherUpper = other.getUpper();
		IPAddress lower = getLower();
		IPAddress upper = getUpper();
		int lowerComp = compareLowValues(lower, otherLower);
		if(!overlaps(other)) {
			if(lowerComp >= 0) {
				if(otherUpper.increment(1).equals(lower)) {
					return create(otherLower, upper);
				}
			} else {
				if(upper.increment(1).equals(otherLower)) {
					return create(lower, otherUpper);
				}
			}
			return null;
		}
		int upperComp = compareLowValues(upper, otherUpper);
		IPAddress lowestLower, highestUpper;
		if(lowerComp >= 0) {
			if(lowerComp == 0 && upperComp == 0) {
				return this;
			}
			lowestLower = otherLower;
		} else {
			lowestLower = lower;
		}
		highestUpper = upperComp >= 0 ? upper : otherUpper;
		return create(lowestLower, highestUpper);
	}

	/**
	 * Extend this sequential range to include all address in the given range, which can be an IPAddress or IPAddressSeqRange.
	 * If the argument has a different IP version than this, null is returned.
	 * Otherwise, this method returns the range that includes this range, the given range, and all addresses in-between.
	 * 
	 * @param other
	 * @return
	 */
	public IPAddressSeqRange extend(IPAddressRange other) {
		IPAddress otherLower = other.getLower();
		IPAddress otherUpper = other.getUpper();
		IPAddress lower = getLower();
		IPAddress upper = getUpper();
		int lowerComp = compareLowValues(lower, otherLower);
		int upperComp = compareLowValues(upper, otherUpper);
		if(lowerComp > 0) { // 
			if(upperComp <= 0) { // ol l u ou
				return other.toSequentialRange();
			}
			IPAddress max = otherUpper.getNetwork().getNetworkMask(getBitCount(), false);
			int versionComp = compareLowValues(lower, max);
			if(versionComp > 0) { // different versions: ol ou max l u
				return null;
			}
			// ol l ou u
			return create(otherLower, upper);
		}
		// lowerComp <= 0
		if(upperComp >= 0) { // l ol ou u
			return this;
		}
		IPAddress max = upper.getNetwork().getNetworkMask(getBitCount(), false);
		int versionComp = compareLowValues(otherLower, max);
		if(versionComp > 0) { // different versions: l u max ol ou
			return null;
		}
		return create(lower, otherUpper);// l ol u ou
	}

	/**
	 * Subtracts the given range from this range, to produce either zero, one, or two address ranges that contain the addresses in this range and not in the given range.
	 * If the result has length 2, the two ranges are ordered by ascending lowest range value. 
	 * 
	 * @param other
	 * @return
	 */
	public IPAddressSeqRange[] subtract(IPAddressSeqRange other) {
		IPAddress otherLower = other.getLower();
		IPAddress otherUpper = other.getUpper();
		IPAddress lower = getLower();
		IPAddress upper = getUpper();
		if(compareLowValues(lower, otherLower) < 0) {
			if(compareLowValues(upper, otherUpper) > 0) { // l ol ou u
				return createPair(lower, otherLower.increment(-1), otherUpper.increment(1), upper);
			} else {
				int comp = compareLowValues(upper, otherLower);
				if(comp < 0) { // l u ol ou
					return createSingle();
				} else if(comp == 0) { // l u == ol ou
					return createSingle(lower, upper.increment(-1));
				}
				return createSingle(lower, otherLower.increment(-1)); // l ol u ou 
			}
		} else if(compareLowValues(otherUpper, upper) >= 0) { // ol l u ou
			return createEmpty();
		} else {
			int comp = compareLowValues(otherUpper, lower);
			if(comp < 0) {
				return createSingle(); // ol ou l u
			} else if(comp == 0) {
				return createSingle(lower.increment(1), upper); // ol ou == l u
			}
			return createSingle(otherUpper.increment(1), upper); // ol l ou u    
		}
	}
	
	protected abstract IPAddressSeqRange create(IPAddress lower, IPAddress upper);
	
	protected abstract IPAddressSeqRange[] createPair(IPAddress lower1, IPAddress upper1, IPAddress lower2, IPAddress upper2);
	
	protected abstract IPAddressSeqRange[] createSingle(IPAddress lower, IPAddress upper);
	
	protected abstract IPAddressSeqRange[] createSingle();
	
	protected abstract IPAddressSeqRange[] createEmpty();

	@Override
	public boolean containsPrefixBlock(int prefixLen) {
		IPAddressSection.checkSubnet(lower, prefixLen);
		int divCount = lower.getDivisionCount();
		int bitsPerSegment = lower.getBitsPerSegment();
		int i = getHostSegmentIndex(prefixLen, lower.getBytesPerSegment(), bitsPerSegment);
		if(i < divCount) {
			IPAddressSegment div = lower.getSegment(i);
			IPAddressSegment upperDiv = upper.getSegment(i);
			int segmentPrefixLength = IPAddressSection.getPrefixedSegmentPrefixLength(bitsPerSegment, prefixLen, i);
			if(!div.containsPrefixBlock(div.getSegmentValue(), upperDiv.getSegmentValue(), segmentPrefixLength)) {
				return false;
			}
			for(++i; i < divCount; i++) {
				div = lower.getSegment(i);
				upperDiv = upper.getSegment(i);
				//is full range?
				if(!div.includesZero() || !upperDiv.includesMax()) {
					return false;
				}
			}
		}
		return true;
	}
	
	@Override
	public boolean containsSinglePrefixBlock(int prefixLen) {
		IPAddressSection.checkSubnet(lower, prefixLen);
		int prevBitCount = 0;
		int divCount = lower.getDivisionCount();
		for(int i = 0; i < divCount; i++) {
			IPAddressSegment div = lower.getSegment(i);
			IPAddressSegment upperDiv = upper.getSegment(i);
			int bitCount = div.getBitCount();
			int totalBitCount = bitCount + prevBitCount;
			if(prefixLen >= totalBitCount) {
				if(!div.isSameValues(upperDiv)) {
					return false;
				}
			} else  {
				int divPrefixLen = Math.max(0, prefixLen - prevBitCount);
				if(!div.containsSinglePrefixBlock(div.getSegmentValue(), upperDiv.getSegmentValue(), divPrefixLen)) {
					return false;
				}
				for(++i; i < divCount; i++) {
					div = lower.getSegment(i);
					upperDiv = upper.getSegment(i);
					//is full range?
					if(!div.includesZero() || !upperDiv.includesMax()) {
						return false;
					}
				}
				return true;
			}
			prevBitCount = totalBitCount;
		}
		return true;
	}
	
	@Override
	public int getBitCount() {
		return getLower().getBitCount();
	}

	@Override
	public byte[] getBytes() {
		return getLower().getBytes();
	}

	@Override
	public byte[] getBytes(byte[] bytes) {
		return getLower().getBytes(bytes);
	}

	@Override
	public byte[] getBytes(byte[] bytes, int index) {
		return getLower().getBytes(bytes, index);
	}

	@Override
	public byte[] getUpperBytes() {
		return getUpper().getUpperBytes();
	}

	@Override
	public byte[] getUpperBytes(byte[] bytes) {
		return getUpper().getUpperBytes(bytes);
	}

	@Override
	public byte[] getUpperBytes(byte[] bytes, int index) {
		return getUpper().getUpperBytes(bytes, index);
	}

	@Override
	public BigInteger getValue() {
		return getLower().getValue();
	}

	@Override
	public BigInteger getUpperValue() {
		return getUpper().getValue();
	}

	@Override
	public boolean isZero() {
		return includesZero() && !isMultiple();
	}

	@Override
	public boolean includesZero() {
		return getLower().isZero();
	}

	@Override
	public boolean isMax() {
		return includesMax() && !isMultiple();
	}

	@Override
	public boolean includesMax() {
		return getUpper().isMax();
	}
}

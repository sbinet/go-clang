package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
import "C"

/**
 * \brief Describes the availability of a particular entity, which indicates
 * whether the use of this entity will result in a warning or error due to
 * it being deprecated or unavailable.
 */
type AvailabilityKind uint32

const (

	/**
	 * \brief The entity is available.
	 */
	Available AvailabilityKind = C.CXAvailability_Available

	/**
	 * \brief The entity is available, but has been deprecated (and its use is
	 * not recommended).
	 */
	Deprecated AvailabilityKind = C.CXAvailability_Deprecated
	/**
	 * \brief The entity is not available; any use of it will be an error.
	 */
	NotAvailable AvailabilityKind = C.CXAvailability_NotAvailable
	/**
	 * \brief The entity is available, but not accessible; any use of it will be
	 * an error.
	 */
	NotAccessible AvailabilityKind = C.CXAvailability_NotAccessible
)

func (ak AvailabilityKind) String() string {
	switch ak {
	case Available:
		return "Available"
	case Deprecated:
		return "Deprecated"
	case NotAvailable:
		return "NotAvailable"
	case NotAccessible:
		return "NotAccessible"
	}
	return ""
}

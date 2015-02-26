package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

/**
 * Describes the availability of a given entity on a particular platform, e.g.,
 * a particular class might only be available on Mac OS 10.7 or newer.
 */
type PlatformAvailability struct {
	c C.CXPlatformAvailability
}

/**
 * \brief A string that describes the platform for which this structure
 * provides availability information.
 *
 * Possible values are "ios" or "macosx".
 */
func (p *PlatformAvailability) Platform() string {
	o := cxstring{p.c.Platform}
	//defer o.Dispose() // done by PlatformAvailability.Dispose()
	return o.String()
}

/**
 * \brief The version number in which this entity was introduced.
 */
func (p *PlatformAvailability) Introduced() Version {
	o := Version{p.c.Introduced}
	return o
}

/**
 * \brief The version number in which this entity was deprecated (but is
 * still available).
 */
func (p *PlatformAvailability) Deprecated() Version {
	o := Version{p.c.Deprecated}
	return o
}

/**
 * \brief The version number in which this entity was obsoleted, and therefore
 * is no longer available.
 */
func (p *PlatformAvailability) Obsoleted() Version {
	o := Version{p.c.Obsoleted}
	return o
}

/**
 * \brief Whether the entity is unconditionally unavailable on this platform.
 */
func (p *PlatformAvailability) Unavailable() int {
	return int(p.c.Unavailable)
}

/**
 * \brief An optional message to provide to a user of this API, e.g., to
 * suggest replacement APIs.
 */
func (p *PlatformAvailability) Message() string {
	o := cxstring{p.c.Message}
	//defer o.Dispose() // done by PlatformAvailability.Dispose()
	return o.String()
}

/**
 * \brief Free the memory associated with a \c CXPlatformAvailability structure.
 */
func (p *PlatformAvailability) Dispose() {
	C.clang_disposeCXPlatformAvailability(&p.c)
}

package clang

// #include "clang-c/Index.h"
import "C"

/**
 * \brief Describes the severity of a particular diagnostic.
 */
type DiagnosticSeverity int

const (
	/**
	 * \brief A diagnostic that has been suppressed, e.g., by a command-line
	 * option.
	 */
	Diagnostic_Ignored = C.CXDiagnostic_Ignored

	/**
	 * \brief This diagnostic is a note that should be attached to the
	 * previous (non-note) diagnostic.
	 */
	Diagnostic_Note = C.CXDiagnostic_Note

	/**
	 * \brief This diagnostic indicates suspicious code that may not be
	 * wrong.
	 */
	Diagnostic_Warning = C.CXDiagnostic_Warning

	/**
	 * \brief This diagnostic indicates that the code is ill-formed.
	 */
	Diagnostic_Error = C.CXDiagnostic_Error

	/**
	 * \brief This diagnostic indicates that the code is ill-formed such
	 * that future parser recovery is unlikely to produce useful
	 * results.
	 */
	Diagnostic_Fatal = C.CXDiagnostic_Fatal
)

func (ds DiagnosticSeverity) String() string {
	switch ds {
	case Diagnostic_Ignored:
		return "Ignored"
	case Diagnostic_Note:
		return "Note"
	case Diagnostic_Warning:
		return "Warning"
	case Diagnostic_Error:
		return "Error"
	case Diagnostic_Fatal:
		return "Fatal"
	default:
		return "Invalid"
	}
}

/**
 * \brief A single diagnostic, containing the diagnostic's severity,
 * location, text, source ranges, and fix-it hints.
 */
type Diagnostic struct {
	c C.CXDiagnostic
}

type Diagnostics []Diagnostic

func (d Diagnostics) Dispose() {
	for _, di := range d {
		di.Dispose()
	}
}

/**
 * \brief Destroy a diagnostic.
 */
func (d Diagnostic) Dispose() {
	C.clang_disposeDiagnostic(d.c)
}

/**
 * \brief Determine the severity of the given diagnostic.
 */
func (d Diagnostic) Severity() DiagnosticSeverity {
	return DiagnosticSeverity(C.clang_getDiagnosticSeverity(d.c))
}

/**
 * \brief Retrieve the source location of the given diagnostic.
 *
 * This location is where Clang would print the caret ('^') when
 * displaying the diagnostic on the command line.
 */
func (d Diagnostic) Location() SourceLocation {
	return SourceLocation{C.clang_getDiagnosticLocation(d.c)}
}

/**
 * \brief Retrieve the text of the given diagnostic.
 */
func (d Diagnostic) Spelling() string {
	cx := cxstring{C.clang_getDiagnosticSpelling(d.c)}
	defer cx.Dispose()
	return cx.String()
}

/**
 * \brief Retrieve the name of the command-line option that enabled this
 * diagnostic.
 *
 * \param Diag The diagnostic to be queried.
 *
 * \param Disable If non-NULL, will be set to the option that disables this
 * diagnostic (if any).
 *
 * \returns A string that contains the command-line option used to enable this
 * warning, such as "-Wconversion" or "-pedantic".
 */
func (d Diagnostic) Option() (enable, disable string) {
	var c_disable cxstring
	cx := cxstring{C.clang_getDiagnosticOption(d.c, &c_disable.c)}
	defer cx.Dispose()
	defer c_disable.Dispose()
	return cx.String(), c_disable.String()
}

/**
 * \brief Retrieve a source range associated with the diagnostic.
 *
 * A diagnostic's source ranges highlight important elements in the source
 * code. On the command line, Clang displays source ranges by
 * underlining them with '~' characters.
 *
 * \param Diagnostic the diagnostic whose range is being extracted.
 *
 * \param Range the zero-based index specifying which range to
 *
 * \returns the requested source range.
 */
func (d Diagnostic) Ranges() (ret []SourceRange) {
	ret = make([]SourceRange, C.clang_getDiagnosticNumRanges(d.c))
	for i := range ret {
		ret[i].c = C.clang_getDiagnosticRange(d.c, C.uint(i))
	}
	return
}

type FixIt struct {
	Data             string
	ReplacementRange SourceRange
}

/**
 * \brief Retrieve the replacement information for a given fix-it.
 *
 * Fix-its are described in terms of a source range whose contents
 * should be replaced by a string. This approach generalizes over
 * three kinds of operations: removal of source code (the range covers
 * the code to be removed and the replacement string is empty),
 * replacement of source code (the range covers the code to be
 * replaced and the replacement string provides the new code), and
 * insertion (both the start and end of the range point at the
 * insertion location, and the replacement string provides the text to
 * insert).
 *
 * \param Diagnostic The diagnostic whose fix-its are being queried.
 *
 * \param FixIt The zero-based index of the fix-it.
 *
 * \param ReplacementRange The source range whose contents will be
 * replaced with the returned replacement string. Note that source
 * ranges are half-open ranges [a, b), so the source code should be
 * replaced from a and up to (but not including) b.
 *
 * \returns A string containing text that should be replace the source
 * code indicated by the \c ReplacementRange.
 */
func (d Diagnostic) FixIts() (ret []FixIt) {
	ret = make([]FixIt, C.clang_getDiagnosticNumFixIts(d.c))
	for i := range ret {
		cx := cxstring{C.clang_getDiagnosticFixIt(d.c, C.uint(i), &ret[i].ReplacementRange.c)}
		defer cx.Dispose()
		ret[i].Data = cx.String()
	}
	return
}

/**
 * \brief Format the given diagnostic in a manner that is suitable for display.
 *
 * This routine will format the given diagnostic to a string, rendering
 * the diagnostic according to the various options given. The
 * \c clang_defaultDiagnosticDisplayOptions() function returns the set of
 * options that most closely mimics the behavior of the clang compiler.
 *
 * \param Diagnostic The diagnostic to print.
 *
 * \param Options A set of options that control the diagnostic display,
 * created by combining \c CXDiagnosticDisplayOptions values.
 *
 * \returns A new string containing for formatted diagnostic.
 */
func (d Diagnostic) Format(options DiagnosticDisplayOptions) string {
	cx := cxstring{C.clang_formatDiagnostic(d.c, C.uint(options))}
	defer cx.Dispose()
	return cx.String()
}

func (d Diagnostic) String() string {
	return d.Format(DefaultDiagnosticDisplayOptions())
}

/**
 * \brief Options to control the display of diagnostics.
 *
 * The values in this enum are meant to be combined to customize the
 * behavior of \c clang_displayDiagnostic().
 */
type DiagnosticDisplayOptions int

const (
	/**
	 * \brief Display the source-location information where the
	 * diagnostic was located.
	 *
	 * When set, diagnostics will be prefixed by the file, line, and
	 * (optionally) column to which the diagnostic refers. For example,
	 *
	 * \code
	 * test.c:28: warning: extra tokens at end of #endif directive
	 * \endcode
	 *
	 * This option corresponds to the clang flag \c -fshow-source-location.
	 */
	Diagnostic_DisplaySourceLocation = C.CXDiagnostic_DisplaySourceLocation

	/**
	 * \brief If displaying the source-location information of the
	 * diagnostic, also include the column number.
	 *
	 * This option corresponds to the clang flag \c -fshow-column.
	 */
	Diagnostic_DisplayColumn = C.CXDiagnostic_DisplayColumn

	/**
	 * \brief If displaying the source-location information of the
	 * diagnostic, also include information about source ranges in a
	 * machine-parsable format.
	 *
	 * This option corresponds to the clang flag
	 * \c -fdiagnostics-print-source-range-info.
	 */
	Diagnostic_DisplaySourceRanges = C.CXDiagnostic_DisplaySourceRanges

	/**
	 * \brief Display the option name associated with this diagnostic, if any.
	 *
	 * The option name displayed (e.g., -Wconversion) will be placed in brackets
	 * after the diagnostic text. This option corresponds to the clang flag
	 * \c -fdiagnostics-show-option.
	 */
	Diagnostic_DisplayOption = C.CXDiagnostic_DisplayOption

	/**
	 * \brief Display the category number associated with this diagnostic, if any.
	 *
	 * The category number is displayed within brackets after the diagnostic text.
	 * This option corresponds to the clang flag
	 * \c -fdiagnostics-show-category=id.
	 */
	Diagnostic_DisplayCategoryId = C.CXDiagnostic_DisplayCategoryId

	/**
	 * \brief Display the category name associated with this diagnostic, if any.
	 *
	 * The category name is displayed within brackets after the diagnostic text.
	 * This option corresponds to the clang flag
	 * \c -fdiagnostics-show-category=name.
	 */
	Diagnostic_DisplayCategoryName = C.CXDiagnostic_DisplayCategoryName
)

/**
 * \brief Retrieve the set of display options most similar to the
 * default behavior of the clang compiler.
 *
 * \returns A set of display options suitable for use with \c
 * clang_displayDiagnostic().
 */
func DefaultDiagnosticDisplayOptions() DiagnosticDisplayOptions {
	return DiagnosticDisplayOptions(C.clang_defaultDiagnosticDisplayOptions())
}

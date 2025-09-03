# Sentinel Documentation

This directory contains comprehensive documentation for the Sentinel project in multiple formats.

## Contents

1. `sentinel_documentation.tex` - Main LaTeX document with complete technical documentation
2. `diagrams.md` - System diagrams in Mermaid format
3. `generate_visualizations.py` - Python script to generate performance and security charts
4. `figures/` - Directory containing generated visualization images

## Compiling the LaTeX Document

To compile the LaTeX document into a PDF, you'll need a LaTeX distribution installed on your system.

### Using pdflatex

```bash
cd /Users/swayamsingal/Desktop/Programming/Sentinel
pdflatex sentinel_documentation.tex
```

For better results with references and table of contents, run multiple times:

```bash
pdflatex sentinel_documentation.tex
pdflatex sentinel_documentation.tex
```

### Using latexmk (recommended)

```bash
cd /Users/swayamsingal/Desktop/Programming/Sentinel
latexmk -pdf sentinel_documentation.tex
```

### Using Overleaf

You can also upload the `.tex` file to Overleaf (https://www.overleaf.com) for online compilation.

## Generated Visualizations

The following visualizations have been generated and are included in the LaTeX document:

1. `figures/performance_metrics.png` - Performance comparison charts
2. `figures/security_effectiveness.png` - Security effectiveness metrics
3. `figures/crypto_performance.png` - Cryptographic components performance
4. `figures/system_architecture.png` - System architecture diagram
5. `figures/security_pipeline.png` - Security pipeline flow
6. `figures/test_results.png` - Test results visualization

## Updating Visualizations

To regenerate the visualizations:

```bash
cd /Users/swayamsingal/Desktop/Programming/Sentinel
python generate_visualizations.py
```

## Converting to Other Formats

### To Microsoft Word (.docx)

Using pandoc:

```bash
pandoc sentinel_documentation.tex -o sentinel_documentation.docx
```

### To HTML

```bash
pandoc sentinel_documentation.tex -o sentinel_documentation.html
```

## Directory Structure

```
.
├── sentinel_documentation.tex     # Main LaTeX document
├── diagrams.md                    # Mermaid diagrams
├── generate_visualizations.py     # Python visualization script
├── figures/                       # Generated visualization images
│   ├── performance_metrics.png
│   ├── security_effectiveness.png
│   ├── crypto_performance.png
│   ├── system_architecture.png
│   ├── security_pipeline.png
│   └── test_results.png
└── DOCUMENTATION_README.md        # This file
```

## Requirements

- LaTeX distribution (TeX Live, MacTeX, or MiKTeX)
- Python 3.x with matplotlib, seaborn, pandas, and numpy
- Optional: pandoc for format conversion

## Document Sections

The LaTeX document includes:

1. Abstract
2. Introduction
3. Problem Statement
4. Related Work
5. System Architecture
6. Core Components
7. Implementation Details
8. Security Algorithms and Protocols
9. Performance Evaluation
10. Future Work
11. Conclusion
12. References

All diagrams and charts are included as high-resolution PNG images for professional publication quality.

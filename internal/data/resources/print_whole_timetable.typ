#let PAGE_HEIGHT = 297mm
#let PAGE_WIDTH = 420mm
#let PAGE_BORDER = (top:15mm, bottom: 12mm, left: 15mm, right: 15mm)
#set page(height: PAGE_HEIGHT, width: PAGE_WIDTH,
//  numbering: "1",
  margin: PAGE_BORDER,
)
#let CELL_BORDER = 1pt
#let BIG_SIZE = 16pt
#let NORMAL_SIZE = 14pt
#let PLAIN_SIZE = 12pt
#let SMALL_SIZE = 10pt

#let FRAME_COLOUR = "#707070"
#let BREAK_COLOUR = "#e0e0e0"
#let EMPTY_COLOUR = "#f0f0f0"

#let TITLE_HEIGHT = 20mm
#let PLAN_AREA_HEIGHT = (PAGE_HEIGHT - PAGE_BORDER.top
    - PAGE_BORDER.bottom - TITLE_HEIGHT)
#let PLAN_AREA_WIDTH = (PAGE_WIDTH - PAGE_BORDER.left
    - PAGE_BORDER.right)

//#PLAN_AREA_WIDTH x #PLAN_AREA_HEIGHT

#let DAYS = ("Mo", "Di", "Mi", "Do", "Fr")
#let HOURS = ("HU 1", "HU 2",
    "FS 1", "FS 2", "FS 3", "FS 4", "FS 5", "FS 6", "FS 7")

#let H_HEADER_HEIGHT1 = 10mm
#let H_HEADER_HEIGHT2 = 10mm
#let H_HEADER_HEIGHT = H_HEADER_HEIGHT1 + H_HEADER_HEIGHT2
#let V_HEADER_WIDTH = 30mm
#let ROW_HEIGHT = 15mm

// Collect headers and y-coordinates for the rows.
#let ROWS = ("Room 1", "Room2", "Room 3", "Room 4")
#let trows = (H_HEADER_HEIGHT1, H_HEADER_HEIGHT2) + (ROW_HEIGHT,)*ROWS.len()
//#trows

// Build the vertical lines
#let vlines = (V_HEADER_WIDTH,)
#let pcols = DAYS.len()*HOURS.len()
#let colwidth = (PLAN_AREA_WIDTH - V_HEADER_WIDTH) / pcols
#let tcolumns = (V_HEADER_WIDTH,) + (colwidth,)*pcols
//#tcolumns


//#let ch = ([],) + DAYS
//#for h in ptime {
//    ch += (h,) + ([],) * DAYS.len()
//}

#show table.cell: it => {
  if it.y < 2 {
    set text(size: PLAIN_SIZE, weight: "bold")
    align(center + horizon, it.body.at("text", default: ""))
  } else if it.x == 0 {
    set text(size: PLAIN_SIZE, weight: "bold")
    align(center + horizon, it.body.at("text", default: ""))
  } else {
    it
  }
}

// On lines with two text items:
// If one is smaller than 25% of the space, leave this and shrink the
// other to 90% of the reamining space. Otherwise shrink both.
#let fit2inspace(width, text1, text2) = {
    let t1 = text(size: SMALL_SIZE, text1)
    let t2 = text(size: SMALL_SIZE, text2)
    let w4 = width / 4
    context {
        let s1 = measure(t1)
        let s2 = measure(t2)
        if (s1.width + s2.width) > width * 0.9 {
            if s1.width < w4 {
                // shrink only text2
                let w2 = width - s1.width
                let scl = (w2 * 0.9) / s2.width
                box(width: width, inset: 2pt,
                    t1
                    + h(1fr)
                    + text(size: scl * SMALL_SIZE, text2)
                )
            } else if s2.width < w4 {
                // shrink only text1
                let w2 = width - s2.width
                let scl = (w2 * 0.9) / s1.width
                box(width: width, inset: 2pt,
                    text(size: scl * SMALL_SIZE, text1)
                    + h(1fr)
                    + t2
                )
            } else {
                // shrink both
                let scl = (width * 0.9) / (s1.width + s2.width)
                box(width: width, inset: 2pt,
                    text(size: scl * SMALL_SIZE, text1)
                    + h(1fr)
                    + text(size: scl * SMALL_SIZE, text2)
                )
            }
        } else {
            box(width: width, inset: 2pt, t1 + h(1fr) + t2)
        }
    }
}

#let fitinspace(width, textc) = {
    let t = text(size: NORMAL_SIZE, weight: "bold", textc)
    context {
        let s = measure(t)
        if s.width > width * 0.9 {
            let scl = (width * 0.9 / s.width)
            let ts = text(size: scl * NORMAL_SIZE, weight: "bold", textc)
            box(width: width, h(1fr) + ts + h(1fr))
        } else {
            box(width: width, h(1fr) + t + h(1fr))
        }
    }
}

#let cell_inset = CELL_BORDER
#let cell_width = colwidth - cell_inset * 2

#let ttxcell(
    row: 0,
    day: 0,
    hour: 0,
    duration: 1,
    offset: 0,
    fraction: 1,
    total: 1,
    centre: "",
    tl: "",
    tr: "",
    bl: "",
    br: "",
) = {
    let x0 = (day * HOURS.len() + hour) * colwidth + V_HEADER_WIDTH
    let y0 = row * ROW_HEIGHT + H_HEADER_HEIGHT
    let x1 = x0 + colwidth * duration
    let wfrac = cell_width * fraction / total
    let xshift = cell_width * offset / total
    // Shrink excessively large components.
    let b = box(
        fill: luma(100%),
        stroke: CELL_BORDER,
        inset: 0pt,
        height: y1 - y0 - CELL_BORDER*2,
        width: wfrac,
    )[
        #fit2inspace(wfrac, tl, tr)
        #v(1fr)
        #fitinspace(wfrac, centre)
        #v(1fr)
        #fit2inspace(wfrac, bl, br)
    ]
    place(top + left,
        dx: x0 + CELL_BORDER + xshift,
        dy: y0 + CELL_BORDER,
        b
    )
}

#let dheader = []
#let pheader = []

#table(
    columns: tcolumns,
    rows: trows,
    gutter: 0pt,
    stroke: rgb(FRAME_COLOUR),
    inset: 0pt,
    fill: (x, y) =>
        if y > 1 and x > 0 {
            rgb(EMPTY_COLOUR)
        } else {
            rgb(BREAK_COLOUR)
        },
//  align: center + horizon,
    table.header(
        [],
        table.cell(colspan: HOURS.len(), [Montag]),
        table.cell(colspan: HOURS.len(), [Dienstag]),
        table.cell(colspan: HOURS.len(), [Mittwoch]),
        table.cell(colspan: HOURS.len(), [Donnerstag]),
        table.cell(colspan: HOURS.len(), [Freitag]),
        [], table.cell(colspan: pcols, []),
    ),
    [One], table.cell(colspan: pcols, []),
    [Two], table.cell(colspan: pcols, []),
    [Three], table.cell(colspan: pcols, []),
)

/*
#show heading: it => text(weight: "bold", size: BIG_SIZE,
    bottom-edge: "descender",
    pad(left: 5mm, it))
*/
/*
#let xdata = json(sys.inputs.ifile)

#let page = 0
#for (k, kdata) in xdata.Pages [
    #{
        if page != 0 {
            pagebreak()
        }
        page += 1
    }

    = #k

    #box([
        #tbody
        #for kd in kdata {
            ttcell(..kd)
        }
    ])
]
*/

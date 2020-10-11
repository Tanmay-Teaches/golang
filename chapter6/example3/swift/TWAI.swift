import Foundation

public struct Tetris {
    typealias Piece = Int

    struct PieceType {
        static let I = 0
        static let J = 1
        static let L = 2
        static let S = 3
        static let Z = 4
        static let T = 5
        static let O = 6

        static let coords: [((Float, Float), (Float, Float), (Float, Float), (Float, Float), (Float, Float))] = [
            ((2, 1), (0.5, 0.5), (1.5, 0.5), (2.5, 0.5), (3.5, 0.5)),
            ((1.5, 0.5), (0.5, 0.5), (1.5, 0.5), (2.5, 0.5), (2.5, -0.5)),
            ((1.5, 0.5), (0.5, 0.5), (1.5, 0.5), (2.5, 0.5), (2.5, 1.5)),
            ((1.5, 1.5), (0.5, 2.5), (0.5, 1.5), (1.5, 1.5), (1.5, 0.5)),
            ((1.5, 1.5), (0.5, 0.5), (0.5, 1.5), (1.5, 1.5), (1.5, 2.5)),
            ((1.5, 1.5), (0.5, 1.5), (1.5, 0.5), (1.5, 1.5), (1.5, 2.5)),
            ((-1, -1), (0.5, 0.5), (0.5, 1.5), (1.5, 1.5), (1.5, 0.5))
        ]
    }

    struct MovingPiece {
        var type: Piece
        var coords: ((Float, Float), (Float, Float), (Float, Float), (Float, Float), (Float, Float))

        init(type: Piece, at location: (Int, Int)) {
            let location = (Float(location.0), Float(location.1))
            self.type = type
            let coord = PieceType.coords[type]
            let origin = (coord.0.0 + location.0, coord.0.1 + location.1)
            let c1 = (coord.1.0 + location.0, coord.1.1 + location.1)
            let c2 = (coord.2.0 + location.0, coord.2.1 + location.1)
            let c3 = (coord.3.0 + location.0, coord.3.1 + location.1)
            let c4 = (coord.4.0 + location.0, coord.4.1 + location.1)
            self.coords = (origin, c1, c2, c3, c4)
        }
    }

    var width: Int
    var height: Int
    var concreteBoard: [[Int]]
    var held: MovingPiece
    var nextPieces: [MovingPiece]
    var score = 0

    init(width: Int, height: Int) {
        self.width = width
        self.height = height
        concreteBoard = [[Int]](repeating: [Int](repeating: 0, count: width), count: height)
        held = MovingPiece(type: Int.random(in: 0...5), at: (0, width / 2 - 1))
        nextPieces = []
        _ = spawn()
    }

    func pieceHitGround() -> Bool {
        let coords = nextPieces[0].coords
        let c1 = coords.1
        let c2 = coords.2
        let c3 = coords.3
        let c4 = coords.4
        return !freeTile(below: c1) || !freeTile(below: c2) || !freeTile(below: c3) || !freeTile(below: c4)
    }

    static func cwRotate(point: (Float, Float), about origin: (Float, Float)) -> (Float, Float) {
        let normalizedPoint = (point.0 - origin.0, point.1 - origin.1)
        let swappedPoint = (normalizedPoint.1, -normalizedPoint.0)
        return (swappedPoint.0 + origin.0, swappedPoint.1 + origin.1)
    }

    static func ccwRotate(point: (Float, Float), about origin: (Float, Float)) -> (Float, Float) {
        var point = point
        var origin = origin
        for _ in 1...3 {
            point = cwRotate(point: point, about: origin)
            origin = cwRotate(point: origin, about: origin)
        }
        return point
    }

    func freeTile(at point: (Float, Float)) -> Bool {
        let iPoint = (Int(point.0), Int(point.1))
        if iPoint.0 >= height || point.0 < 0 {
            return false
        }
        if iPoint.1 >= width || point.1 < 0 {
            return false
        }
        return concreteBoard[iPoint.0][iPoint.1] == 0
    }

    func freeTile(below point: (Float, Float)) -> Bool {
        if (point.0 - 0.5) == Float(height - 1) {
            return false
        }
        return concreteBoard[Int(point.0) + 1][Int(point.1)] == 0
    }

    mutating func attemptSpin(clockwise: Bool) -> Bool {
        let piece = nextPieces[0]
        if piece.type == PieceType.O {
            return true
        }
        let coords = piece.coords
        var origin = coords.0
        var c1 = coords.1
        var c2 = coords.2
        var c3 = coords.3
        var c4 = coords.4
        let rotator = clockwise ? Tetris.cwRotate : Tetris.ccwRotate
        c1 = rotator(c1, origin)
        c2 = rotator(c2, origin)
        c3 = rotator(c3, origin)
        c4 = rotator(c4, origin)
        origin = rotator(origin, origin)
        if !freeTile(at: c1) ||
           !freeTile(at: c2) ||
           !freeTile(at: c3) ||
           !freeTile(at: c4) {
            return false
        }
        nextPieces[0].coords = (origin, c1, c2, c3, c4)
        return true
    }

    mutating func spawn() -> Bool {
        let nextPiece = MovingPiece(type: Int.random(in: 0...5), at: (0, width / 2 - 1))
        let c1 = nextPiece.coords.1
        let c2 = nextPiece.coords.2
        let c3 = nextPiece.coords.3
        let c4 = nextPiece.coords.4
        if freeTile(at: c1) && freeTile(at: c2) && freeTile(at: c3) && freeTile(at: c4) {
            nextPieces.append(nextPiece)
            return true
        }
        return false
    }

    mutating func lock() -> Bool {
        let coords = nextPieces[0].coords
        let c1 = coords.1
        let c2 = coords.2
        let c3 = coords.3
        let c4 = coords.4
        concreteBoard[Int(c1.0)][Int(c1.1)] = 1
        concreteBoard[Int(c2.0)][Int(c2.1)] = 1
        concreteBoard[Int(c3.0)][Int(c3.1)] = 1
        concreteBoard[Int(c4.0)][Int(c4.1)] = 1
        nextPieces.removeFirst()
        clearLines()
        return spawn()
    }

    mutating func down() -> Bool {
        let coords = nextPieces[0].coords
        let origin = (coords.0.0 + 1, coords.0.1)
        let c1 = (coords.1.0 + 1, coords.1.1)
        let c2 = (coords.2.0 + 1, coords.2.1)
        let c3 = (coords.3.0 + 1, coords.3.1)
        let c4 = (coords.4.0 + 1, coords.4.1)
        nextPieces[0].coords = (origin, c1, c2, c3, c4)
        return pieceHitGround()
    }

    mutating func horizontalMove(left: Bool) -> Bool {
        let piece = nextPieces[0]
        let coords = piece.coords
        let origin = (coords.0.0, coords.0.1 + (left ? -1 : 1))
        let c1 = (coords.1.0, coords.1.1 + (left ? -1 : 1))
        let c2 = (coords.2.0, coords.2.1 + (left ? -1 : 1))
        let c3 = (coords.3.0, coords.3.1 + (left ? -1 : 1))
        let c4 = (coords.4.0, coords.4.1 + (left ? -1 : 1))
        if !freeTile(at: c1) ||
           !freeTile(at: c2) ||
           !freeTile(at: c3) ||
           !freeTile(at: c4) {
            return false
        }
        nextPieces[0].coords = (origin, c1, c2, c3, c4)
        return true
    }

    mutating func swapHold() {
        let newPiece = held
        held = MovingPiece(type: nextPieces.removeFirst().type, at: (0, width / 2 - 1))
        nextPieces.insert(newPiece, at: 0)
    }

    static func clearLines(on board: [[Int]], height: Int, width: Int) -> ([[Int]], Int) {
        var newBoard = board
        var cleared: [Int] = []
        cleared.reserveCapacity(height)
        rowLoop: for i in 0..<height {
            for j in 0..<width {
                if newBoard[i][j] == 0 {
                    continue rowLoop
                }
            }

            cleared.insert(i, at: 0)
        }
        for i in cleared {
            newBoard.remove(at: i)
        }
        while newBoard.count != height {
            newBoard.insert([Int](repeating: 0, count: width), at: 0)
        }
        return (newBoard, cleared.count)
    }

    mutating func clearLines() {
        let (a, cleared) = Tetris.clearLines(on: concreteBoard, height: height, width: width)
        concreteBoard = a
        score += cleared
    }

    func render() -> [[Int]] {
        var board = concreteBoard
        let coords = nextPieces[0].coords
        let c1 = coords.1
        let c2 = coords.2
        let c3 = coords.3
        let c4 = coords.4
        board[Int(c1.0)][Int(c1.1)] = 1
        board[Int(c2.0)][Int(c2.1)] = 1
        board[Int(c3.0)][Int(c3.1)] = 1
        board[Int(c4.0)][Int(c4.1)] = 1
        return board
    }
}

// Weights chosen by CMA-ES
// Reached a top score of over 400,000
public struct TetrisHeuristicWeights: KeyPathIterable {
    var holeCountWeight: Float = 61.156956
    var openHoleCountWeight: Float = 98.64193
    var blocksAboveHolesWeight: Float = 4.9402747
    var nonTetrisClearPenalty: Float = -8.505343
    var tetrisRewardWeight: Float = 10.133134
    var maximumLineHeightWeight: Float = -2.1128924
    var lastBlockAddedHeightWeight: Float = 15.88824
    var pillarCountWeight: Float = 16.198757
    var blocksInRightmostLaneWeight: Float = -21.910303
    var bumpinessWeight: Float = 15.483239
}

public struct TetrisState {
    var height: Int
    var width: Int
    var board: [[Int]]
    var lastBlockAddedHeight: Int

    let weights: TetrisHeuristicWeights

    var holeCount = 0
    var openHoleCount = 0
    var blocksAboveHoles = 0
    var pillarCount = 0
    var maximumLineHeight = 0
    var blocksInRightmostLane = 0
    var bumpiness = 0

    mutating func countHoles() {
        for i in 0..<width {
            var blockFound = false
            var numberOfBlocksFound = 0

            for j in 0..<height {
                if board[j][i] != 0 {
                    blockFound = true
                    numberOfBlocksFound += 1
                } else if blockFound {
                    blocksAboveHoles += numberOfBlocksFound

                    if i < width - 2 {
                        if board[j][i + 1] == 0 && board[j][i + 2] == 0 {
                            if j == height - 1 || board[j + 1][i + 1] != 0 {
                                openHoleCount += 1
                                continue
                            }
                        }
                    }

                    if i >= 2 {
                        if board[j][i - 1] == 0 && board[j][i - 2] == 0 {
                            if j == height - 1 || board[j + 1][i - 1] != 0 {
                                openHoleCount += 1
                                continue
                            }
                        }
                    }

                    holeCount += 1
                }
            }
        }
    }

    mutating func countPillars() {
        for i in 0..<width {
            var currentPillarHeightL = 0
            var currentPillarHeightR = 0

            for j in (0..<height).reversed() {
                if i > 0 && board[j][i] != 0 && board[j][i - 1] == 0 {
                    currentPillarHeightL += 1
                } else {
                    if currentPillarHeightL >= 3 {
                        pillarCount += currentPillarHeightL
                    }
                    currentPillarHeightL = 0
                }

                if i < width - 2 && board[j][i] != 0 && board[j][i + 1] == 0 {
                    currentPillarHeightR += 1
                } else {
                    if currentPillarHeightR >= 3 {
                        pillarCount += currentPillarHeightR
                    }
                    currentPillarHeightR = 0
                }
            }

            if currentPillarHeightL >= 3 {
                pillarCount += currentPillarHeightL
            }
            if currentPillarHeightR >= 3 {
                pillarCount += currentPillarHeightR
            }
        }
    }

    mutating func calculateMaximumLineHeight() {
        for i in 0..<width {
            for j in 0..<height {
                if board[j][i] != 0 {
                    maximumLineHeight = max(maximumLineHeight, height - j)
                    break
                }
            }
        }
    }

    mutating func countBlocksInRightmostLane() {
        for j in 0..<height {
            if board[j][width - 1] != 0 {
                blocksInRightmostLane += 1
            }
        }
    }

    mutating func calculateBumpiness() {
        var previousLineHeight = 0
        for i in 0..<width - 1 {
            for j in 0..<height {
                if board[j][i] != 0 || j == height - 1 {
                    let currentLineHeight = height - j
                    if i != 0 {
                        bumpiness += abs(previousLineHeight - currentLineHeight)
                    }
                    previousLineHeight = currentLineHeight
                    break
                }
            }
        }
    }

    mutating func cost() -> Float {
        let (clearedBoard, linesCleared) = Tetris.clearLines(on: board, height: height, width: width)
        self.board = clearedBoard
        let linesClearedWithoutTetrises = linesCleared > 0 && linesCleared < 4 ? 1 : 0
        let tetrises = linesCleared == 4 ? 1 : 0

        countHoles()
        countPillars()
        calculateMaximumLineHeight()
        countBlocksInRightmostLane()
        calculateBumpiness()

        let cost =
            Float(holeCount) * weights.holeCountWeight +
            Float(openHoleCount) * weights.openHoleCountWeight +
            Float(blocksAboveHoles) * weights.blocksAboveHolesWeight +
            Float(linesClearedWithoutTetrises) * weights.nonTetrisClearPenalty +
            Float(tetrises) * weights.tetrisRewardWeight +
            Float(maximumLineHeight) * weights.maximumLineHeightWeight +
            Float(pillarCount) * weights.pillarCountWeight +
            Float(blocksInRightmostLane) * weights.blocksInRightmostLaneWeight +
            Float(bumpiness) * weights.bumpinessWeight +
            Float(lastBlockAddedHeight) * weights.lastBlockAddedHeightWeight

        return cost
    }
}

public struct TetrisAI {
    struct PathfinderCoordinateGroup: Equatable, Hashable {
        struct Coordinate: Equatable, Hashable {
            var row: Float
            var col: Float

            var tuple: (Float, Float) {
                (row, col)
            }

            init(_ x: (Float, Float)) {
                self.row = x.0
                self.col = x.1
            }

            func isValid(on board: [[Int]]) -> Bool {
                if Int(row) >= board.count || Int(col) >= board[0].count {
                    return false
                }
                if row < 0 || col < 0 {
                    return false
                }
                return board[Int(row)][Int(col)] == 0
            }

            static func ==(lhs: Coordinate, rhs: Coordinate) -> Bool {
                return lhs.row == rhs.row && lhs.col == rhs.col
            }
        }

        var origin: Coordinate
        var c1: Coordinate
        var c2: Coordinate
        var c3: Coordinate
        var c4: Coordinate

        func attemptSpin(on board: [[Int]], clockwise: Bool) -> PathfinderCoordinateGroup? {
            let rotator = clockwise ? Tetris.cwRotate : Tetris.ccwRotate
            let c1 = Coordinate(rotator((self.c1.row, self.c1.col), (origin.row, origin.col)))
            let c2 = Coordinate(rotator((self.c2.row, self.c2.col), (origin.row, origin.col)))
            let c3 = Coordinate(rotator((self.c3.row, self.c3.col), (origin.row, origin.col)))
            let c4 = Coordinate(rotator((self.c4.row, self.c4.col), (origin.row, origin.col)))
            let newOrigin = Coordinate(rotator((origin.row, origin.col), (origin.row, origin.col)))
            if c1.isValid(on: board) && c2.isValid(on: board) && c3.isValid(on: board) && c4.isValid(on: board) {
                return PathfinderCoordinateGroup(origin: newOrigin, c1: c1, c2: c2, c3: c3, c4: c4)
            }
            return nil
        }

        func attemptHorizontal(on board: [[Int]], direction: Float) -> PathfinderCoordinateGroup? {
            let newOrigin = Coordinate((origin.row, origin.col + direction))
            let c1 = Coordinate((self.c1.row, self.c1.col + direction))
            let c2 = Coordinate((self.c2.row, self.c2.col + direction))
            let c3 = Coordinate((self.c3.row, self.c3.col + direction))
            let c4 = Coordinate((self.c4.row, self.c4.col + direction))
            if c1.isValid(on: board) && c2.isValid(on: board) && c3.isValid(on: board) && c4.isValid(on: board) {
                return PathfinderCoordinateGroup(origin: newOrigin, c1: c1, c2: c2, c3: c3, c4: c4)
            }
            return nil
        }

        func attemptDown(on board: [[Int]]) -> PathfinderCoordinateGroup? {
            let newOrigin = Coordinate((origin.row + 1, origin.col))
            let c1 = Coordinate((self.c1.row + 1, self.c1.col))
            let c2 = Coordinate((self.c2.row + 1, self.c2.col))
            let c3 = Coordinate((self.c3.row + 1, self.c3.col))
            let c4 = Coordinate((self.c4.row + 1, self.c4.col))
            if c1.isValid(on: board) && c2.isValid(on: board) && c3.isValid(on: board) && c4.isValid(on: board) {
                return PathfinderCoordinateGroup(origin: newOrigin, c1: c1, c2: c2, c3: c3, c4: c4)
            }
            return nil
        }

        static func ==(lhs: PathfinderCoordinateGroup, rhs: PathfinderCoordinateGroup) -> Bool {
            return lhs.origin == rhs.origin && lhs.c1 == rhs.c1 && lhs.c2 == rhs.c2 && lhs.c3 == rhs.c3 && lhs.c4 == rhs.c4
        }
    }

    static func potentiallyValidFinalStates(in game: Tetris) -> [(PathfinderCoordinateGroup, Int)] {
        var possible: [(PathfinderCoordinateGroup, Int)] = []
        let piece = game.nextPieces[0]
        for row in 0..<game.height {
            let row = Float(row)
            for col in 0..<game.width {
                rotLoop: for rot in 0...3 {
                    let col = Float(col)
                    let coords = Tetris.PieceType.coords[piece.type]
                    var origin = coords.0
                    var c1 = coords.1
                    var c2 = coords.2
                    var c3 = coords.3
                    var c4 = coords.4
                    for _ in 0..<rot {
                        c1 = Tetris.cwRotate(point: c1, about: origin)
                        c2 = Tetris.cwRotate(point: c2, about: origin)
                        c3 = Tetris.cwRotate(point: c3, about: origin)
                        c4 = Tetris.cwRotate(point: c4, about: origin)
                        origin = Tetris.cwRotate(point: origin, about: origin)
                    }
                    origin = (origin.0 + row, origin.1 + col)
                    c1 = (c1.0 + row, c1.1 + col)
                    c2 = (c2.0 + row, c2.1 + col)
                    c3 = (c3.0 + row, c3.1 + col)
                    c4 = (c4.0 + row, c4.1 + col)
                    let cg = PathfinderCoordinateGroup(coords: (origin, c1, c2, c3, c4))
                    if !cg.c1.isValid(on: game.concreteBoard) || !cg.c2.isValid(on: game.concreteBoard) || !cg.c3.isValid(on: game.concreteBoard) || !cg.c4.isValid(on: game.concreteBoard) {
                        continue
                    }
                    if !game.freeTile(at: c1) ||
                        !game.freeTile(at: c2) ||
                        !game.freeTile(at: c3) ||
                        !game.freeTile(at: c4) {
                        continue
                    }
                    if game.freeTile(below: c1) &&
                        game.freeTile(below: c2) &&
                        game.freeTile(below: c3) &&
                        game.freeTile(below: c4) {
                        continue
                    }
                    let highestY = min(min(min(c1.0, c2.0), c3.0), c4.0)
                    possible.append((PathfinderCoordinateGroup(coords: (origin, c1, c2, c3, c4)), game.height - Int(highestY)))
                }
            }
        }
        return possible
    }

    static func findMovesToFinalState(board: [[Int]], start: PathfinderCoordinateGroup, goal: PathfinderCoordinateGroup) -> [Int]? {
        var q: [PathfinderCoordinateGroup] = [start]
        var parents: [PathfinderCoordinateGroup: (PathfinderCoordinateGroup, Int)] = [:]
        parents.reserveCapacity(24 * 10)
        var edges: [PathfinderCoordinateGroup] = []
        var edgeActions: [Int] = []
        for _ in 1...7 {
            edges.append(.init(coords: ((0, 0), (0, 0), (0, 0), (0, 0), (0, 0))))
            edgeActions.append(-1)
        }
        while !q.isEmpty {
            let v = q.removeFirst()
            if v == goal {
                var actions: [Int] = []
                var lastEdge = v
                while lastEdge != start {
                    let parent = parents[lastEdge]!
                    actions.append(parent.1)
                    lastEdge = parent.0
                }
                return actions.reversed()
            }
            var edgeCount = 0
            if let clockwise = v.attemptSpin(on: board, clockwise: true) {
                edges[edgeCount] = clockwise
                edgeActions[edgeCount] = 0
                edgeCount += 1
            }
            if let counterclockwise = v.attemptSpin(on: board, clockwise: false) {
                edges[edgeCount] = counterclockwise
                edgeActions[edgeCount] = 2
                edgeCount += 1
            }
            if let right = v.attemptHorizontal(on: board, direction: 1) {
                edges[edgeCount] = right
                edgeActions[edgeCount] = 4
                edgeCount += 1
            }
            if let left = v.attemptHorizontal(on: board, direction: -1) {
                edges[edgeCount] = left
                edgeActions[edgeCount] = 5
                edgeCount += 1
            }
            if let down = v.attemptDown(on: board) {
                edges[edgeCount] = down
                edgeActions[edgeCount] = 6
                edgeCount += 1
            }
            for i in 0..<edgeCount {
                let edge = edges[i]
                if !parents.keys.contains(edge) {
                    parents[edge] = (v, edgeActions[i])
                    q.append(edge)
                }
            }
        }
        return nil
    }
}

extension TetrisAI.PathfinderCoordinateGroup {
    init(coords: ((Float, Float), (Float, Float), (Float, Float), (Float, Float), (Float, Float))) {
        origin = Coordinate(coords.0)
        c1 = .init(coords.1)
        c2 = .init(coords.2)
        c3 = .init(coords.3)
        c4 = .init(coords.4)
    }
}

extension Tetris {
    mutating func execute(moves: [Int]) {
        for i in moves {
            switch i {
            case -1:
                swapHold()
            case 0:
                _ = attemptSpin(clockwise: true)
            case 2:
                _ = attemptSpin(clockwise: false)
            case 4:
                _ = horizontalMove(left: false)
            case 5:
                _ = horizontalMove(left: true)
            case 6:
                _ = down()
            default:
                fatalError()
            }
        }
    }

    func executeOnCopies(moves: [Int]) -> [Tetris] {
        var copies: [Tetris] = [self]
        for i in moves {
            copies.append(copies.last!)
            switch i {
            case -1:
                copies[copies.count - 1].swapHold()
            case 0:
                _ = copies[copies.count - 1].attemptSpin(clockwise: true)
            case 2:
                _ = copies[copies.count - 1].attemptSpin(clockwise: false)
            case 4:
                _ = copies[copies.count - 1].horizontalMove(left: false)
            case 5:
                _ = copies[copies.count - 1].horizontalMove(left: true)
            case 6:
                _ = copies[copies.count - 1].down()
            default:
                fatalError()
            }
        }
        return copies
    }
}

extension Tetris {
    func nextBestMoves(weights: TetrisHeuristicWeights = TetrisHeuristicWeights(), lookAtHold: Bool = true) -> ([Int], Float)? {
        var potentialFinalStates: [(TetrisAI.PathfinderCoordinateGroup, Float, Int)] = TetrisAI.potentiallyValidFinalStates(in: self).map { ($0.0, -1, $0.1) }
        var board = concreteBoard
        for i in 0..<potentialFinalStates.count {
            let currentMove = potentialFinalStates[i].0
            board[Int(currentMove.c1.row)][Int(currentMove.c1.col)] = 1
            board[Int(currentMove.c2.row)][Int(currentMove.c2.col)] = 1
            board[Int(currentMove.c3.row)][Int(currentMove.c3.col)] = 1
            board[Int(currentMove.c4.row)][Int(currentMove.c4.col)] = 1
            var state = TetrisState(height: height, width: width, board: board, lastBlockAddedHeight: potentialFinalStates[i].2, weights: weights)
            potentialFinalStates[i].1 = state.cost()
            board[Int(currentMove.c1.row)][Int(currentMove.c1.col)] = 0
            board[Int(currentMove.c2.row)][Int(currentMove.c2.col)] = 0
            board[Int(currentMove.c3.row)][Int(currentMove.c3.col)] = 0
            board[Int(currentMove.c4.row)][Int(currentMove.c4.col)] = 0
        }
        let sortedStates = potentialFinalStates.sorted(by: { $0.1 < $1.1 })
        var potentialHoldMoves: ([Int], Float)? = nil
        if lookAtHold {
            var copy = self
            copy.swapHold()
            potentialHoldMoves = copy.nextBestMoves(weights: weights, lookAtHold: false)
        }
        for i in sortedStates {
            if let moves = TetrisAI.findMovesToFinalState(board: concreteBoard, start: TetrisAI.PathfinderCoordinateGroup(coords: nextPieces[0].coords), goal: i.0) {
                if lookAtHold {
                    if let holdMoves = potentialHoldMoves {
                        if holdMoves.1 < i.1 {
                            return ([-1] + holdMoves.0, holdMoves.1)
                        }
                    }
                }
                return (moves, i.1)
            }
        }
        return nil
    }
}

public struct SimulatedGame {
    static let gamesPerIndividual = 2

    static func playGame(with weights: TetrisHeuristicWeights) -> Float {
        var scores: [Int] = []
        scores.reserveCapacity(gamesPerIndividual)
        for _ in 1...gamesPerIndividual {
            var game = Tetris(width: 10, height: 24)
            while true {
                if let moves = game.nextBestMoves(weights: weights) {
                    game.execute(moves: moves.0)
                }
                if !game.lock() {
                    scores.append(game.score)
                    break
                }
            }
        }
        return Float(scores.reduce(0, +)) / Float(gamesPerIndividual)
    }
}

var game = Tetris(width: 10, height: 24)

@_cdecl("nextBestMoves")
public func nextBestMoves() -> UnsafeMutablePointer<Int> {
    var nextMoves = game.nextBestMoves()!.0
    nextMoves.insert(nextMoves.count, at: 0)
    return nextMoves.withUnsafeBufferPointer { ptrToMoves -> UnsafeMutablePointer<Int> in
        let newMemory = UnsafeMutablePointer<Int>.allocate(capacity: nextMoves.count)
        memcpy(newMemory, nextMoves, nextMoves.count * MemoryLayout<Int>.size)
        return newMemory
    }
}

@_cdecl("playMove")
public func playMove(move: Int) {
    switch move {
    case -1:
        game.swapHold()
    case 0:
        game.attemptSpin(clockwise: true)
    case 2:
        game.attemptSpin(clockwise: false)
    case 4:
        game.horizontalMove(left: false)
    case 5:
        game.horizontalMove(left: true)
    case 6:
        game.down()
    default:
        fatalError()
    }
}
@_cdecl("renderFrame")
public func renderFrame() -> UnsafeMutablePointer<Int> {
    var x = [24, 10] + game.render().reduce([], +)
    x.insert(x.count, at: 0)
    return x.withUnsafeBufferPointer { ptrToMoves -> UnsafeMutablePointer<Int> in
        let newMemory = UnsafeMutablePointer<Int>.allocate(capacity: x.count)
        memcpy(newMemory, x, x.count * MemoryLayout<Int>.size)
        return newMemory
    }
}

@_cdecl("lockGame")
public func lockGame() -> Bool {
    return game.lock()
}

@_cdecl("resetGame")
public func resetGame() {
    game = Tetris(width: 10, height: 24)
}
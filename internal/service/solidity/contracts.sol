// SPDX-License-Identifier: GPL-3.0
pragma solidity 0.8.0;

import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/access/Ownable.sol";

interface AggregatorInterface {
  function latestAnswer() external view returns (int256);
  function latestTimestamp() external view returns (uint256);
  function latestRound() external view returns (uint256);
  function getAnswer(uint256 roundId) external view returns (int256);
  function getTimestamp(uint256 roundId) external view returns (uint256);

  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);
  event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt);
}

interface AggregatorV3Interface {

  function decimals() external view returns (uint8);
  function description() external view returns (string memory);
  function version() external view returns (uint256);

  // getRoundData and latestRoundData should both raise "No data present"
  // if they do not have data to report, instead of returning unset values
  // which could be misinterpreted as actual reported values.
  function getRoundData(uint80 _roundId)
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );
  function latestRoundData()
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );

}

interface AggregatorV2V3Interface is AggregatorInterface, AggregatorV3Interface
{
}

contract SimpleAggregator is AggregatorV2V3Interface, Ownable {
    struct Round {
        int256 answer;
        uint64 startedAt;
        uint64 updatedAt;
        uint80 answeredInRound;
    }

    uint80 constant private ROUND_MAX = 2**80-1;
    string constant private V3_NO_DATA_ERROR = "No data present";

    uint8 public override decimals;
    string public override description;
    uint256 constant public override version = 3;
    
    uint80 internal latestRoundId;
    mapping(uint80 => Round) internal rounds;

    constructor(
        uint8 _decimals,
        string memory _description
    ) {
        decimals = _decimals;
        description = _description;
    }

    function set(
        uint80 _roundId,
        int256 _answer
    ) public onlyOwner {
        rounds[_roundId].answer = _answer;
        rounds[_roundId].startedAt = uint64(block.timestamp);
        rounds[_roundId].updatedAt = uint64(block.timestamp);
        rounds[_roundId].answeredInRound = _roundId;
        latestRoundId = _roundId;
    }

    /**
   * @notice get the most recently reported answer
   *
   * @dev #[deprecated] Use latestRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended latestRoundData
   * instead which includes better verification information.
   */
  function latestAnswer()
    public
    view
    virtual
    override
    returns (int256)
  {
    return rounds[latestRoundId].answer;
  }

  /**
   * @notice get the most recent updated at timestamp
   *
   * @dev #[deprecated] Use latestRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended latestRoundData
   * instead which includes better verification information.
   */
  function latestTimestamp()
    public
    view
    virtual
    override
    returns (uint256)
  {
    return rounds[latestRoundId].updatedAt;
  }

  /**
   * @notice get the ID of the last updated round
   *
   * @dev #[deprecated] Use latestRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended latestRoundData
   * instead which includes better verification information.
   */
  function latestRound()
    public
    view
    virtual
    override
    returns (uint256)
  {
    return latestRoundId;
  }

  /**
   * @notice get past rounds answers
   * @param _roundId the round number to retrieve the answer for
   *
   * @dev #[deprecated] Use getRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended getRoundData
   * instead which includes better verification information.
   */
  function getAnswer(uint256 _roundId)
    public
    view
    virtual
    override
    returns (int256)
  {
    if (validRoundId(_roundId)) {
      return rounds[uint80(_roundId)].answer;
    }
    return 0;
  }

  /**
   * @notice get timestamp when an answer was last updated
   * @param _roundId the round number to retrieve the updated timestamp for
   *
   * @dev #[deprecated] Use getRoundData instead. This does not error if no
   * answer has been reached, it will simply return 0. Either wait to point to
   * an already answered Aggregator or use the recommended getRoundData
   * instead which includes better verification information.
   */
  function getTimestamp(uint256 _roundId)
    public
    view
    virtual
    override
    returns (uint256)
  {
    if (validRoundId(_roundId)) {
      return rounds[uint80(_roundId)].updatedAt;
    }
    return 0;
  }

    /**
   * @notice get data about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @param _roundId the round ID to retrieve the round data for
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started. This is 0
   * if the round hasn't been started yet.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed. answeredInRound may be smaller than roundId when the round
   * timed out. answeredInRound is equal to roundId when the round didn't time out
   * and was completed regularly.
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet received
   * maxSubmissions) answer and updatedAt may change between queries.
   */
  function getRoundData(uint80 _roundId)
    public
    view
    virtual
    override
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    require(validRoundId(_roundId), V3_NO_DATA_ERROR);

    Round memory r = rounds[_roundId];

    return (
      _roundId,
      r.answer,
      r.startedAt,
      r.updatedAt,
      r.answeredInRound
    );
  }

  /**
   * @notice get data about the latest round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values. Consumers are encouraged to
   * use this more fully featured method over the "legacy" latestRound/
   * latestAnswer/latestTimestamp functions. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started. This is 0
   * if the round hasn't been started yet.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed. answeredInRound may be smaller than roundId when the round
   * timed out. answeredInRound is equal to roundId when the round didn't time
   * out and was completed regularly.
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet
   * received maxSubmissions) answer and updatedAt may change between queries.
   */
   function latestRoundData()
    public
    view
    virtual
    override
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return getRoundData(latestRoundId);
  }

  function validRoundId(uint256 _roundId)
    private
    pure
    returns (bool)
  {
    return _roundId <= ROUND_MAX;
  }
}

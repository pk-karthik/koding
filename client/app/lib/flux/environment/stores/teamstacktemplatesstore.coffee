KodingFluxStore      = require 'app/flux/base/store'
toImmutable          = require 'app/util/toImmutable'
immutable            = require 'immutable'
actions              = require '../actiontypes'


module.exports = class TeamStackTemplatesStore extends KodingFluxStore

  @getterPath = 'TeamStackTemplatesStore'

  getInitialState: -> immutable.Map()

  initialize: ->

    @on actions.LOAD_TEAM_STACK_TEMPLATES_SUCCESS, @load
    @on actions.CHANGE_TEMPLATE_TITLE, @changeTitle

    @on actions.REMOVE_STACK_TEMPLATE_SUCCESS, @remove
    @on actions.UPDATE_TEAM_STACK_TEMPLATE_SUCCESS, @updateSingle
    @on actions.UPDATE_STACK_TEMPLATE_SUCCESS, @updateSingle


  load: (stackTemplates, { templates }) ->

    stackTemplates = stackTemplates.withMutations (_templates) ->
      templates.forEach (template) ->
        _templates.set template._id, toImmutable template

    return stackTemplates


  changeTitle: (stackTemplates, { id, value }) ->

    template = stackTemplates.get id

    return stackTemplates  unless template

    stackTemplates.withMutations (templates) ->
      templates
        .setIn [id, 'title'], value
        .setIn [id, 'isDirty'], yes


  remove: (stackTemplates, { id }) -> stackTemplates.remove id

  updateSingle: (stackTemplates, { stackTemplate }) ->

    return stackTemplates  if stackTemplate.accessLevel isnt 'group'

    return stackTemplates.set stackTemplate._id, toImmutable stackTemplate
